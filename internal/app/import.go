package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"

	"pickcampus-backend/internal/bootstrap"
	"pickcampus-backend/internal/config"
	"pickcampus-backend/models"
	"pickcampus-backend/models/factory"
)

var importDataDir string

// importCmd 一次性把前端 public/{admission,admission-major} JSON 导入 tbl_admission。
// 用法: pickcampus-backend import-admission --data-dir /path/to/pickcampus/public -c configs/config.yaml
var importCmd = &cobra.Command{
	Use:   "import-admission",
	Short: "从前端 public/{admission,admission-major} JSON 导入录取数据到 tbl_admission",
	RunE: func(*cobra.Command, []string) error {
		return runImport()
	},
}

func init() {
	importCmd.Flags().StringVar(&importDataDir, "data-dir", "",
		"前端 public 目录路径(含 admission/ 与 admission-major/)")
	rootCmd.AddCommand(importCmd)
}

// admissionDTO 前端 JSON 记录(camelCase);major 缺省(nil) = 院校级记录。
type admissionDTO struct {
	UniversityID string             `json:"universityId"`
	Province     string             `json:"province"`
	Subject      string             `json:"subject"`
	Year         int                `json:"year"`
	MinRank      int                `json:"minRank"`
	MinScore     *int               `json:"minScore"`
	Source       string             `json:"source"`
	Major        *string            `json:"major"`
	MajorCode    string             `json:"majorCode"`
	ElectiveReq  string             `json:"electiveReq"`
	Batch        string             `json:"batch"`
	Tuition      *int               `json:"tuition"`
	Duration     string             `json:"duration"`
	Rating       string             `json:"rating"`
	RatingRank   *int               `json:"ratingRank"`
	SubRatings   []models.SubRating `json:"subRatings"`
}

func (d admissionDTO) toModel() *models.Admission {
	return &models.Admission{
		UniversityID: d.UniversityID,
		Province:     d.Province,
		Subject:      d.Subject,
		Year:         d.Year,
		MinRank:      d.MinRank,
		MinScore:     d.MinScore,
		Source:       d.Source,
		Major:        d.Major,
		MajorCode:    d.MajorCode,
		ElectiveReq:  d.ElectiveReq,
		Batch:        d.Batch,
		Tuition:      d.Tuition,
		Duration:     d.Duration,
		Rating:       d.Rating,
		RatingRank:   d.RatingRank,
		SubRatings:   models.SubRatingList(d.SubRatings),
	}
}

func runImport() error {
	if importDataDir == "" {
		return fmt.Errorf("必须用 --data-dir 指定前端 public 目录")
	}
	if err := cleanenv.ReadConfig(configFile, config.G); err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}
	if err := bootstrap.InitDB(config.G.Base.MySQL); err != nil {
		return err
	}
	if err := bootstrap.AutoMigrate(models.AllTables...); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}
	db := bootstrap.Cli(context.Background())
	admissionRepo := factory.AdmissionRepo(db)

	// 先清空重导(幂等,便于反复运行)
	if err := admissionRepo.DeleteAll(); err != nil {
		return fmt.Errorf("清空旧数据失败: %w", err)
	}

	// 流式:逐文件解析 + 插入,峰值内存只占单个省份文件,适配小内存机器
	total := 0
	for _, sub := range []string{"admission", "admission-major"} {
		dir := filepath.Join(importDataDir, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("读取目录 %s 失败: %w", dir, err)
		}
		subTotal := 0
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			recs, err := parseAdmissionFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return err
			}
			if err := admissionRepo.BulkInsert(recs); err != nil {
				return fmt.Errorf("插入 %s 失败: %w", e.Name(), err)
			}
			subTotal += len(recs)
		}
		fmt.Printf("导入 %s: %d 条\n", sub, subTotal)
		total += subTotal
	}

	college, _ := admissionRepo.CountByLevel(false)
	major, _ := admissionRepo.CountByLevel(true)
	fmt.Printf("导入完成:院校级 %d 条,专业级 %d 条,合计 %d 条\n", college, major, college+major)
	return nil
}

// parseAdmissionFile 解析单个 JSON 文件为 Admission 记录(流式导入的最小单元)。
func parseAdmissionFile(path string) ([]*models.Admission, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dtos []admissionDTO
	if err := json.Unmarshal(raw, &dtos); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %w", filepath.Base(path), err)
	}
	out := make([]*models.Admission, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toModel())
	}
	return out, nil
}
