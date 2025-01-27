package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zbednarke/geomag/pkg/egm96"
	"github.com/zbednarke/geomag/pkg/wmm"
)

const (
	usage          = "wmm_point --cof_file=WMM2020.COF --spherical [latitude] [longitude] [altitude] [date]"
	cofUsage       = "COF coefficients file to use, empty for the built-in one"
	sphericalUsage = "Output spherical values instead of ellipsoidal"
	lngErr         = "Error: Degree input is outside legal range. The legal range is from -180 to 360."
	fieldWarn      = "Warning: The Horizontal Field strength at this location is only 0.000000. " +
		"Compass readings have VERY LARGE uncertainties in areas where where H is smaller than 1000 nT"
)

var prompt = map[string]string{
	"latitude": "Please enter latitude North Latitude positive. " +
		"For example: 30, 30, 30 (D,M,S) or 30.508 (Decimal Degrees) (both are north). ",
	"longitude": "Please enter longitude East longitude positive, West negative. " +
		"For example: -100.5 or -100, 30, 0 for 100.5 degrees west. ",
	"altitude": "Please enter height above mean sea level (in kilometers). " +
		"[For height above WGS-84 Ellipsoid prefix E, for example (E20.1)]. ",
	"date": "Please enter the decimal year or calendar date (YYYY.yyy, MM DD YYYY or MM/DD/YYYY) ",
}

type Dataset struct {
	Latitude  []float64 `json:"latitude"`
	Longitude []float64 `json:"longitude"`
	Altitude  []float64 `json:"altitude"`
	Bx        []float64 `json:"bx"`
	By        []float64 `json:"by"`
	Bz        []float64 `json:"bz"`
	// Dbx       []float64 `json:"dbx"`
	// Dby       []float64 `json:"dby"`
	// Dbz       []float64 `json:"dbz"`
}

func (ds *Dataset) WriteToJSON(filename string) error {
	// Create or truncate file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Optionally set indentation for readability
	encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ")

	// Encode the struct to JSON and write directly to the file
	if err := encoder.Encode(ds); err != nil {
		return err
	}

	return nil
}

var (
	cofFile    string
	spherical  bool
	latitude   float64
	longitude  float64
	altitude   float64
	hae        bool
	dYear      float64
	ErrHelp    error
	err        error
	loc        egm96.Location
	x, y, z    float64
	dx, dy, dz float64
	dataset    Dataset
)

func init() {
	flag.StringVar(&cofFile, "cof_file", "", cofUsage)
	flag.StringVar(&cofFile, "c", "", cofUsage)

	flag.BoolVar(&spherical, "spherical", false, sphericalUsage)
	flag.BoolVar(&spherical, "s", false, sphericalUsage)

	ErrHelp = errors.New(usage)
}

func main() {
	flag.Parse()

	// if cofFile != "" {
	// 	if err = wmm.LoadWMMCOF(cofFile); err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }
	fmt.Printf("COF File: %v, Epoch: %v, Valid Date: %d/%d/%d\n", wmm.COFName, wmm.Epoch,
		wmm.ValidDate.Month(), wmm.ValidDate.Day(), wmm.ValidDate.Year())

	// if flag.NArg() == 0 {
	// 	userInput()
	// } else if flag.NArg() == 4 {
	// 	if latitude, err = parsing.ParseLatLng(flag.Arg(0)); err != nil {
	// 		_, _ = fmt.Fprintln(os.Stderr, err)
	// 		return
	// 	}
	// 	if longitude, err = parsing.ParseLatLng(flag.Arg(1)); err != nil {
	// 		_, _ = fmt.Fprintln(os.Stderr, err)
	// 		return
	// 	}
	// 	if altitude, hae, err = parsing.ParseAltitude(flag.Arg(2)); err != nil {
	// 		_, _ = fmt.Fprintln(os.Stderr, err)
	// 		return
	// 	}
	// 	if dYear, err = parsing.ParseTime(flag.Arg(3)); err != nil {
	// 		_, _ = fmt.Fprintln(os.Stderr, err)
	// 		return
	// 	}
	// } else {
	// 	_, _ = fmt.Fprintf(os.Stderr, "You must specify a latitude, longitude, altitude and date in that order")
	// 	return
	// }

	latitude = 0.0
	// longitude = -100.0
	altitude = 0.0

	idx := -1
	for longitude = 0; longitude < 360; longitude += .001 {
		idx += 1
		for longitude < 0 {
			longitude += 360
		}
		if longitude >= 360 {
			_, _ = fmt.Fprintln(os.Stderr, lngErr)
		}
		altitude *= 1000 // Convert to meters

		if hae {
			loc = egm96.NewLocationGeodetic(latitude, longitude, altitude)
		} else {
			loc, err = egm96.NewLocationMSL(latitude, longitude, altitude)
			if err != nil {
				fmt.Printf("Error making location: %s\n", err)
			}
		}
		mf, err := wmm.CalculateWMMMagneticField(
			loc,
			wmm.DecimalYear(dYear).ToTime(),
		)

		if err != nil {
			fmt.Printf("Error calculating magnetic field: %s\n", err)
		}

		// x,y,z,dx,dy,dz, lat, long, alt

		// fmt.Println("Results For")
		// fmt.Println()
		lat, lng, hh := loc.Geodetic()
		// qualifier := "N"
		// quantity := lat / egm96.Deg
		// if quantity < 0 {
		// qualifier = "S"
		// quantity = -quantity
		// }

		// fmt.Printf("Latitude:\t%4.2f%s\n", quantity, qualifier)

		// qualifier = "E"
		// quantity = lng / egm96.Deg
		// if quantity >= 180 {
		// qualifier = "W"
		// quantity = 360 - quantity
		// }
		// fmt.Printf("Longitude:\t%4.2f%s\n", quantity, qualifier)

		// relationship := "above"
		// quantity = hh
		// qualifier = "the WGS-84 ellipsoid"
		// if !hae {
		// quantity, _ = loc.HeightAboveMSL()
		// qualifier = "mean sea level"
		// }
		// if quantity < 0 {
		// relationship = "below"
		// quantity = -quantity
		// }
		// fmt.Printf("Altitude:\t%6.3f kilometers %s %s\n", quantity/1000, relationship, qualifier)

		// fmt.Printf("Date:\t\t%5.1f\n", dYear)

		// qualifier = ""
		// if spherical {
		// 	qualifier = "(Spherical)"
		// }
		// fmt.Println()

		// if err != nil {
		// fmt.Printf("Warning: %s\n\n", err)
		// }

		if spherical {
			x, y, z, dx, dy, dz = mf.Spherical()
		} else {
			x, y, z, dx, dy, dz = mf.Ellipsoidal()
		}

		dataset.Latitude = append(dataset.Latitude, lat)
		dataset.Longitude = append(dataset.Longitude, lng)
		dataset.Altitude = append(dataset.Altitude, hh)
		dataset.Bx = append(dataset.Bx, x)
		dataset.By = append(dataset.By, y)
		dataset.Bz = append(dataset.Bz, z)
		// dataset.Dbx = append(dataset.Dbx, dx)
		// dataset.Dby = append(dataset.Dby, dy)
		// dataset.Dbz = append(dataset.Dbz, dz)

		// dD, dM, dS := egm96.DegreesToDMS(mf.D())
		// iD, iM, iS := egm96.DegreesToDMS(mf.I())
		// gvD, gvM, gvS := egm96.DegreesToDMS(mf.GV(loc))
		// fmt.Println("       Main Field             Secular Change")
		// fmt.Printf("F    = %8.1f nT ± %5.1f nT  %6.1f nT/yr\n", mf.F(), mf.ErrF(), mf.DF())
		// if !spherical {
		// fmt.Printf("H    = %8.1f nT ± %5.1f nT  %6.1f nT/yr\n", mf.H(), mf.ErrH(), mf.DH())
		// }
		// fmt.Printf("X    = %8.1f nT ± %5.1f nT  %6.1f nT/yr %s\n", x, mf.ErrX(), dx, qualifier)
		// fmt.Printf("Y    = %8.1f nT ± %5.1f nT  %6.1f nT/yr %s\n", y, mf.ErrY(), dy, qualifier)
		// fmt.Printf("Z    = %8.1f nT ± %5.1f nT  %6.1f nT/yr %s\n", z, mf.ErrZ(), dz, qualifier)
		// if !spherical {
		// fmt.Printf("Decl =    %3.0fº %2.0f' ± %2.0f'         %4.1f'/yr\n", dD, dM+dS/60, mf.ErrD()*60, mf.DD()*60)
		// fmt.Printf("Incl =    %3.0fº %2.0f' ± %2.0f'         %4.1f'/yr\n", iD, iM+iS/60, mf.ErrI()*60, mf.DI()*60)
		// fmt.Println()
		// fmt.Printf("Grid Variation =  %2.0fº %2.0f'\n", gvD, gvM+gvS/60)
		// }
	}

	if err := dataset.WriteToJSON("equatorDataset.json"); err != nil {
		panic(err)
	}
}

func readUserInput(prompt string) (inp string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	inp, _ = reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	return inp
}
