package reducers

import (
	"testing"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/shopspring/decimal"
)

func getDecimalFromString(str string) decimal.Decimal {
	decimalVal, _ := decimal.NewFromString(str)
	return decimalVal
}

func TestGetMetadataOfMetric(t *testing.T) {

	metric := common.Metric{
		TimeGrain:      "month",
		Period:         "2023-02",
		Dimension:      "none",
		DimensionValue: "No dimension",
		History: common.MetricHistory{
			"1f38fc06b5af2e642e12855619dd5347e5feb07e": {
				Lines:           2,
				KPI:             getDecimalFromString("110891.3"),
				CommitTimestamp: 1681911827,
			},
			"1f5aebaf23cfe91630ca6671cf6de639b98bc32d": {
				Lines:           2,
				KPI:             getDecimalFromString("111039.3"),
				CommitTimestamp: 1682574676,
			},
			"22413c25649f318c69efde0bfd1aa9ea6e970bee": {
				Lines:           2,
				KPI:             getDecimalFromString("110417.3"),
				CommitTimestamp: 1679548972,
			},
			"2cc6c5e78a1b1cbfdb938b8aa23f24540092fe4e": {
				Lines:           2,
				KPI:             getDecimalFromString("110543.3"),
				CommitTimestamp: 1677667509,
			},
			"364afa96cbdcb098b14a4335e2fd6848200948ee": {
				Lines:           2,
				KPI:             getDecimalFromString("110417.3"),
				CommitTimestamp: 1680585827,
			},
			"43582f8727c32d23a43832c9cc949d75341290d9": {
				Lines:           2,
				KPI:             getDecimalFromString("110891.3"),
				CommitTimestamp: 1682416519,
			},
			"448a392e3cbd96b75e940acbb52a9a91d016904f": {
				Lines:           2,
				KPI:             getDecimalFromString("110348.3"),
				CommitTimestamp: 1679117511,
			},
			"5b3906ec17f955b69660aba158a9e40ebb9ed2a4": {
				Lines:           2,
				KPI:             getDecimalFromString("110417.3"),
				CommitTimestamp: 1679376019,
			},
			"7f0fccb9af6d9ea6ac4547ebb434e8942c2ee817": {
				Lines:           2,
				KPI:             getDecimalFromString("111039.3"),
				CommitTimestamp: 1683006553,
			},
			"7f9bc2762b14a785fc3f215e4ce20bae7802a343": {
				Lines:           2,
				KPI:             getDecimalFromString("110891.3"),
				CommitTimestamp: 1682492113,
			},
			"8590e724ed18ab779f6fce19b9705d5a401337f2": {
				Lines:           2,
				KPI:             getDecimalFromString("110417.3"),
				CommitTimestamp: 1679475234,
			},
			"8be46462b20b29f133893a530a83c60db6e6735c": {
				Lines:           2,
				KPI:             getDecimalFromString("110417.3"),
				CommitTimestamp: 1679053203,
			},
			"9c84209dfe2351df0a05435e8e8c90f9198462a9": {
				Lines:           2,
				KPI:             getDecimalFromString("110891.3"),
				CommitTimestamp: 1681895066,
			},
			"acc6c4d29a81ee8ebb73ca067c8055d2dc085cef": {
				Lines:           2,
				KPI:             getDecimalFromString("111039.3"),
				CommitTimestamp: 1683177691,
			},
			"dc7358e8bb5c00b717d8e835af563c79a3aba049": {
				Lines:           2,
				KPI:             getDecimalFromString("110891.3"),
				CommitTimestamp: 1681917529,
			},
			"e01fed18fc731b0b16d42c221d9f6dde6500bf5c": {
				Lines:           2,
				KPI:             getDecimalFromString("110348.3"),
				CommitTimestamp: 1678985820,
			},
		},
	}

	relativeHistory := map[time.Duration]RelativeHistoricalEvent{
		10*time.Hour + 45*time.Minute + 10*time.Second:   {getDecimalFromString("100").Sub(decimal.NewFromInt(100)), getDecimalFromString("0.4480324074074074"), 1677667509},
		376*time.Hour + 57*time.Minute + 1*time.Second:   {getDecimalFromString("99.82359853559646").Sub(decimal.NewFromInt(100)), getDecimalFromString("15.706261574074075"), 1678985820},
		395*time.Hour + 40*time.Minute + 4*time.Second:   {getDecimalFromString("99.88601751530848").Sub(decimal.NewFromInt(100)), getDecimalFromString("16.486157407407408"), 1679053203},
		413*time.Hour + 31*time.Minute + 52*time.Second:  {getDecimalFromString("99.82359853559646").Sub(decimal.NewFromInt(100)), getDecimalFromString("17.230462962962964"), 1679117511},
		485*time.Hour + 20*time.Minute + 20*time.Second:  {getDecimalFromString("99.88601751530848").Sub(decimal.NewFromInt(100)), getDecimalFromString("20.222453703703703"), 1679376019},
		512*time.Hour + 53*time.Minute + 55*time.Second:  {getDecimalFromString("99.88601751530848").Sub(decimal.NewFromInt(100)), getDecimalFromString("21.370775462962964"), 1679475234},
		533*time.Hour + 22*time.Minute + 53*time.Second:  {getDecimalFromString("99.88601751530848").Sub(decimal.NewFromInt(100)), getDecimalFromString("22.224224537037035"), 1679548972},
		821*time.Hour + 23*time.Minute + 48*time.Second:  {getDecimalFromString("99.88601751530848").Sub(decimal.NewFromInt(100)), getDecimalFromString("34.22486111111111"), 1680585827},
		1185*time.Hour + 4*time.Minute + 27*time.Second:  {getDecimalFromString("100.31480876724324").Sub(decimal.NewFromInt(100)), getDecimalFromString("49.37809027777778"), 1681895066},
		1189*time.Hour + 43*time.Minute + 48*time.Second: {getDecimalFromString("100.31480876724324").Sub(decimal.NewFromInt(100)), getDecimalFromString("49.57208333333333"), 1681911827},
		1191*time.Hour + 18*time.Minute + 50*time.Second: {getDecimalFromString("100.31480876724324").Sub(decimal.NewFromInt(100)), getDecimalFromString("49.638078703703705"), 1681917529},
		1329*time.Hour + 55*time.Minute + 20*time.Second: {getDecimalFromString("100.31480876724324").Sub(decimal.NewFromInt(100)), getDecimalFromString("55.41342592592593"), 1682416519},
		1350*time.Hour + 55*time.Minute + 14*time.Second: {getDecimalFromString("100.31480876724324").Sub(decimal.NewFromInt(100)), getDecimalFromString("56.288356481481486"), 1682492113},
		1373*time.Hour + 51*time.Minute + 17*time.Second: {getDecimalFromString("100.44869295561106").Sub(decimal.NewFromInt(100)), getDecimalFromString("57.243946759259266"), 1682574676},
		1493*time.Hour + 49*time.Minute + 14*time.Second: {getDecimalFromString("100.44869295561106").Sub(decimal.NewFromInt(100)), getDecimalFromString("62.242523148148145"), 1683006553},
		1541*time.Hour + 21*time.Minute + 32*time.Second: {getDecimalFromString("100.44869295561106").Sub(decimal.NewFromInt(100)), getDecimalFromString("64.22328703703704"), 1683177691},
	}
	firstDate, _ := GetFirstDateOfPeriod("2023-02")
	expected := MetricMetadata{
		TimeGrain:       "month",
		PeriodKey:       "2023-02",
		InitialValue:    getDecimalFromString("110543.3"),
		FirstDate:       firstDate,
		RelativeHistory: relativeHistory,
	}

	result, _ := GetMetadataOfMetric(metric)

	if expected.FirstDate != result.FirstDate {
		t.Errorf("Expected firstDate %s, but got %s", expected.FirstDate, result.FirstDate)
	}

	if !expected.InitialValue.Equal(result.InitialValue) {
		t.Errorf("Expected InitialValue %s, but got %s", expected.InitialValue, result.InitialValue)
	}

	for key, value := range result.RelativeHistory {

		if !expected.RelativeHistory[key].RelativeValue.Equal(value.RelativeValue) {
			t.Errorf("Expected RelativeHistory %s for key %s, but got %s", expected.RelativeHistory[key].RelativeValue, key, value.RelativeValue)
		}

		if expected.RelativeHistory[key].ComputationTimetamp != value.ComputationTimetamp {
			t.Errorf("Expected RelativeHistory %d for key %s, but got %d", expected.RelativeHistory[key].ComputationTimetamp, key, value.ComputationTimetamp)
		}
	}

}
