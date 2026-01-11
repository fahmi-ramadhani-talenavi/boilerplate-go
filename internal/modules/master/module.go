package master

import (
	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/modules/master/handler"
	"github.com/user/go-boilerplate/pkg/cache"
	"gorm.io/gorm"
)

// Module represents the master data module.
type Module struct {
	areaHandler           *handler.AreaHandler
	provinceHandler       *handler.ProvinceHandler
	districtHandler       *handler.DistrictHandler
	bankHandler           *handler.BankHandler
	branchHandler         *handler.BranchHandler
	genderHandler         *handler.GenderHandler
	religionHandler       *handler.ReligionHandler
	maritalStatusHandler  *handler.MaritalStatusHandler
	citizenshipHandler    *handler.CitizenshipHandler
	educationLevelHandler *handler.EducationLevelHandler
	currencyHandler       *handler.CurrencyHandler
	taxGroupHandler       *handler.TaxGroupHandler
	taxBracketHandler     *handler.TaxBracketHandler
	batchHandler         *handler.BatchHandler
}

// New creates a new master module.
func New(db *gorm.DB, cfg *config.Config, cache *cache.Client) *Module {
	return &Module{
		areaHandler:           handler.NewAreaHandler(db),
		provinceHandler:       handler.NewProvinceHandler(db),
		districtHandler:       handler.NewDistrictHandler(db),
		bankHandler:           handler.NewBankHandler(db),
		branchHandler:         handler.NewBranchHandler(db),
		genderHandler:         handler.NewGenderHandler(db),
		religionHandler:       handler.NewReligionHandler(db),
		maritalStatusHandler:  handler.NewMaritalStatusHandler(db),
		citizenshipHandler:    handler.NewCitizenshipHandler(db),
		educationLevelHandler: handler.NewEducationLevelHandler(db),
		currencyHandler:       handler.NewCurrencyHandler(db),
		taxGroupHandler:       handler.NewTaxGroupHandler(db),
		taxBracketHandler:     handler.NewTaxBracketHandler(db),
		batchHandler:         handler.NewBatchHandler(db, cache),
	}
}

// RegisterRoutes registers master data routes.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	master := api.Group("/master")
	master.GET("/all", m.batchHandler.All)
	master.GET("/areas", m.areaHandler.List)
	master.GET("/provinces", m.provinceHandler.List)
	master.GET("/districts", m.districtHandler.List)
	master.GET("/banks", m.bankHandler.List)
	master.GET("/branches", m.branchHandler.List)
	master.GET("/genders", m.genderHandler.List)
	master.GET("/religions", m.religionHandler.List)
	master.GET("/marital-statuses", m.maritalStatusHandler.List)
	master.GET("/citizenships", m.citizenshipHandler.List)
	master.GET("/education-levels", m.educationLevelHandler.List)
	master.GET("/currencies", m.currencyHandler.List)
	master.GET("/tax-groups", m.taxGroupHandler.List)
	master.GET("/tax-brackets", m.taxBracketHandler.List)
}
