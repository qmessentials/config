package utilities

import (
	"config/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func BootstrapConfig(db *gorm.DB) error {
	var bootstrapped bool
	temp := db.Raw(`select "values"[1]::boolean as is_bootstrapped from config_settings where name = 'IsBootstrapped'`)
	if temp.Error != nil {
		return temp.Error
	}
	if err := temp.Scan(&bootstrapped).Error; err != nil {
		return err
	}
	if bootstrapped {
		return nil
	}
	log.Warn().Msg("Bootstrapping config info")
	units := []models.Unit{
		{FullName: "inch", FullNamePlural: "inches", Abbreviation: "in", MeasurementSystem: "US", UnitType: "linear"},
		{FullName: "foot", FullNamePlural: "feet", Abbreviation: "ft", MeasurementSystem: "US", UnitType: "linear"},
		{FullName: "yard", FullNamePlural: "yards", Abbreviation: "yd", MeasurementSystem: "US", UnitType: "linear"},
		{FullName: "mile", FullNamePlural: "miles", Abbreviation: "mi", MeasurementSystem: "US", UnitType: "linear"},
		{FullName: "square inch", FullNamePlural: "square inches", Abbreviation: "sq. in.", MeasurementSystem: "US", UnitType: "area"},
		{FullName: "square foot", FullNamePlural: "square feet", Abbreviation: "sq ft", MeasurementSystem: "US", UnitType: "area"},
		{FullName: "square yard", FullNamePlural: "square yards", Abbreviation: "sq yd", MeasurementSystem: "US", UnitType: "area"},
		{FullName: "acre", FullNamePlural: "acres", Abbreviation: "ac", MeasurementSystem: "US", UnitType: "area"},
		{FullName: "square mile", FullNamePlural: "square miles", Abbreviation: "sq mi", MeasurementSystem: "US", UnitType: "area"},
		{FullName: "cubic inch", FullNamePlural: "cubic inches", Abbreviation: "cu in", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "cubic foot", FullNamePlural: "cubic feet", Abbreviation: "cu ft", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "cubic yard", FullNamePlural: "cubic yards", Abbreviation: "cu yd", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "teaspoon", FullNamePlural: "teaspoons", Abbreviation: "tsp", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "tablespoon", FullNamePlural: "tablespoons", Abbreviation: "tbsp", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "cup", FullNamePlural: "cups", Abbreviation: "c", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "pint", FullNamePlural: "pints", Abbreviation: "pt", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "quart", FullNamePlural: "quarts", Abbreviation: "qt", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "gallon", FullNamePlural: "gallons", Abbreviation: "gal", MeasurementSystem: "US", UnitType: "volume"},
		{FullName: "ounce", FullNamePlural: "ounces", Abbreviation: "oz", MeasurementSystem: "US", UnitType: "weight"},
		{FullName: "pound", FullNamePlural: "pounds", Abbreviation: "lb", MeasurementSystem: "US", UnitType: "weight"},
		{FullName: "ton", FullNamePlural: "tons", Abbreviation: "ton", MeasurementSystem: "US", UnitType: "weight"},
		{FullName: "mile per hour", FullNamePlural: "miles per hour", Abbreviation: "mph", MeasurementSystem: "US", UnitType: "velocity"},
		{FullName: "foot per second", FullNamePlural: "feet per second", Abbreviation: "ft/s", MeasurementSystem: "US", UnitType: "velocity"},
		{FullName: "yard per second", FullNamePlural: "yards per second", Abbreviation: "yd/s", MeasurementSystem: "US", UnitType: "velocity"},
		{FullName: "inch per second squared", FullNamePlural: "inches per second squared", Abbreviation: "in/s^2", MeasurementSystem: "US", UnitType: "acceleration"},
		{FullName: "foot per second squared", FullNamePlural: "feet per second squared", Abbreviation: "ft/s^2", MeasurementSystem: "US", UnitType: "acceleration"},
		{FullName: "yard per second squared", FullNamePlural: "yards per second squared", Abbreviation: "yd/s^2", MeasurementSystem: "US", UnitType: "acceleration"},
		{FullName: "pound per square inch", FullNamePlural: "pounds per square inch", Abbreviation: "psi", MeasurementSystem: "US", UnitType: "pressure"},
		{FullName: "pound per square foot", FullNamePlural: "pounds per square foot", Abbreviation: "psf", MeasurementSystem: "US", UnitType: "pressure"},
		{FullName: "pound per square yard", FullNamePlural: "pounds per square yard", Abbreviation: "psy", MeasurementSystem: "US", UnitType: "pressure"},
		{FullName: "millimeter", FullNamePlural: "millimeters", Abbreviation: "mm", MeasurementSystem: "metric", UnitType: "linear"},
		{FullName: "centimeter", FullNamePlural: "centimeters", Abbreviation: "cm", MeasurementSystem: "metric", UnitType: "linear"},
		{FullName: "meter", FullNamePlural: "meters", Abbreviation: "m", MeasurementSystem: "metric", UnitType: "linear"},
		{FullName: "kilometer", FullNamePlural: "kilometers", Abbreviation: "km", MeasurementSystem: "metric", UnitType: "linear"},
		{FullName: "square millimeter", FullNamePlural: "square millimeters", Abbreviation: "sq mm", MeasurementSystem: "metric", UnitType: "area"},
		{FullName: "square centimeter", FullNamePlural: "square centimeters", Abbreviation: "sq cm", MeasurementSystem: "metric", UnitType: "area"},
		{FullName: "square meter", FullNamePlural: "square meters", Abbreviation: "sq m", MeasurementSystem: "metric", UnitType: "area"},
		{FullName: "hectare", FullNamePlural: "hectares", Abbreviation: "ha", MeasurementSystem: "metric", UnitType: "area"},
		{FullName: "square kilometer", FullNamePlural: "square kilometers", Abbreviation: "sq km", MeasurementSystem: "metric", UnitType: "area"},
		{FullName: "cubic millimeter", FullNamePlural: "cubic millimeters", Abbreviation: "cu mm", MeasurementSystem: "metric", UnitType: "volume"},
		{FullName: "cubic centimeter", FullNamePlural: "cubic centimeters", Abbreviation: "cu cm", MeasurementSystem: "metric", UnitType: "volume"},
		{FullName: "cubic meter", FullNamePlural: "cubic meters", Abbreviation: "cu m", MeasurementSystem: "metric", UnitType: "volume"},
		{FullName: "milliliter", FullNamePlural: "milliliters", Abbreviation: "mL", MeasurementSystem: "metric", UnitType: "volume"},
		{FullName: "liter", FullNamePlural: "liters", Abbreviation: "L", MeasurementSystem: "metric", UnitType: "volume"},
		{FullName: "gram", FullNamePlural: "grams", Abbreviation: "g", MeasurementSystem: "metric", UnitType: "mass"},
		{FullName: "kilogram", FullNamePlural: "kilograms", Abbreviation: "kg", MeasurementSystem: "metric", UnitType: "mass"},
		{FullName: "metric ton", FullNamePlural: "metric tons", Abbreviation: "t", MeasurementSystem: "metric", UnitType: "mass"},
		{FullName: "meter per second", FullNamePlural: "meters per second", Abbreviation: "m/s", MeasurementSystem: "metric", UnitType: "velocity"},
		{FullName: "kilometer per hour", FullNamePlural: "kilometers per hour", Abbreviation: "km/h", MeasurementSystem: "metric", UnitType: "velocity"},
		{FullName: "meter per second squared", FullNamePlural: "meters per second squared", Abbreviation: "m/s^2", MeasurementSystem: "metric", UnitType: "acceleration"},
		{FullName: "kilogram per square meter", FullNamePlural: "kilograms per square meter", Abbreviation: "kg/m^2", MeasurementSystem: "metric", UnitType: "pressure"},
		{FullName: "second", FullNamePlural: "seconds", Abbreviation: "s", MeasurementSystem: "none", UnitType: "time"},
		{FullName: "minute", FullNamePlural: "minutes", Abbreviation: "min", MeasurementSystem: "none", UnitType: "time"},
		{FullName: "hour", FullNamePlural: "hours", Abbreviation: "hr", MeasurementSystem: "none", UnitType: "time"},
		{FullName: "day", FullNamePlural: "days", Abbreviation: "day", MeasurementSystem: "none", UnitType: "time"},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		//Apparently there's a bug with how GORM handles array inserts in Postgresql
		if err := tx.Exec("INSERT INTO config_settings (name,values) VALUES ('IsBootstrapped',ARRAY['true'])").Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO config_settings (name,values) VALUES ('Modifiers',ARRAY['top','bottom','left','right','middle','upper','lower','center','inside','outside','warp','fill'])").Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO config_settings (name,values) VALUES ('MeasurementSystems',ARRAY['metric','US','none'])").Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO config_settings (name,values) VALUES ('UnitTypes',ARRAY['linear','area','volume','weight','mass','velocity','acceleration','pressure','time'])").Error; err != nil {
			return err
		}
		if err := db.Create(&units).Error; err != nil {
			return err
		}
		return nil
	})
	return err

}
