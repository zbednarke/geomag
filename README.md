# geomag
geomag is an implementation in Go of the NOAA World Magnetic Model.

The World Magnetic Model home is at https://www.ngdc.noaa.gov/geomag/WMM/DoDWMM.shtml.

The coefficients for 2020-2024 can be downloaded at https://www.ngdc.noaa.gov/geomag/WMM/data/WMM2020/WMM2020COF.zip

## Commands
geomag provides two command line programs, modeled after the command line programs in the official NOAA software.

`wmm_point` calculates magnetic field values for a single location and time:
```
> wmm_point N30 W88.51 0.01 2019.5

Results For

Latitude:       30.00N
Longitude:      88.51W
Altitude:        0.010 kilometers above mean sea level
Date:           2019.5

       Main Field             Secular Change
       F    =  46944.3 nT ± 152.0 nT  -118.8 nT/yr
       H    =  24074.6 nT ± 133.0 nT    -6.8 nT/yr
       X    =  24060.2 nT ± 138.0 nT    -8.0 nT/yr
       Y    =   -831.0 nT ±  89.0 nT   -36.3 nT/yr
       Z    =  40301.2 nT ± 165.0 nT  -134.3 nT/yr
       Decl =     -1º 59' ± 19'         -5.2'/yr
       Incl =     59º  9' ± 13'         -4.6'/yr

       Grid Variation =  -1º 59'
```

`wmm_grid` is coming soon.  It will calculate magnetic field values for a grid of locations and/or times.

## Packages
Two packages are provided by this library:

### egm96
Package egm96 provides a representation of the 1996 Earth Gravitational Model (EGM96),
a geopotential model of the Earth.

EGM96 is the geoid reference model component of the World Geodetic System (WGS84).
It consists of n=m=360 spherical harmonic coefficients as published by the
National Geospatial-Intelligence Agency (NGA).  The NGA also publishes a raster grid
of the calculated heights which can be interpolated to approximate the geoid height
at any location.

In effect, this model provides the height of sea level above the WGS84 reference ellipsoid.
It is used, for example, in GPS navigation to provide the height above sea level.

This package is based on the NGA-provided 15'x15' resolution grid encoding
the heights of the geopotential surface at each lat/long, and interpolates between grid
points using a bilinear interpolation.

usage:
```
import "github.com/zbednarke/geomag/pkg/egm96"

// Calculate height above MSL for a point at a
// latitude of 12.25 South, longitude of 82.75 East, and
// altitude of 1000m above the WGS84 ellipsoid (i.e. GPS altitude)
h, err := egm96.NewLocationGeodetic(-12.25, 82.75, 1000).HeightAboveMSL()
```

### wmm
Package wmm provides a representation of the 2020 World Magnetic Model (WMM),
a mathematical model of the magnetic field produced by the Earth's core and
its variation over time.

WMM is the magnetic model component of the World Geodetic System (WGS84).
It consists of n=m=12 spherical harmonic coefficients as published by the
National Geospatial-Intelligence Agency (NGA).

This model evaluates all magnetic field components and their rates of change
for any location on the Earth's surface.  These field components include the
X, Y, and Z values and the declination D and inclination I.
The Declination is used, for example, in correcting a Magnetic Heading to a
True Heading.

usage:
```
import "github.com/zbednarke/geomag/pkg/egm96"
import "github.com/zbednarke/geomag/pkg/wmm"

tt := wmm.DecimalYear(2019.5)
loc := egm96.NewLocationGeodetic(-12.25, 82.75, 1000)

mag, _ := wmm.CalculateWMMMagneticField(loc, tt.ToTime())
fmt.Printf("Declination at your location: %2.2f\n", mag.D())
```

## Validation
The library code is fully tested.
In particular, all test values provided with the official NOAA WMM are tested here,
as well as the detailed example in the WMM technical paper.
Please submit an issue on github if you notice any other issues.

## Updating Coefficients
Use go-bindata in the root directory to update the coefficients stored in binary form.
Coefficients are currently updated through the WMM2020 model.
First, unzip the new WMM zip file in the assets/wmm directory, then
`go-bindata -o ../../../pkg/wmm/bindata.go WMM.COF`
Inside the `bindata.go` file, change the package from `main` to `wmm`.
Change the `Assets` function to `getAssets` and remove the
`MustAsset`, `AssetInfo`, `AssetNames`, `AssetDir`, `bintree`, `_bintree`, `RestoreAsset`, `RestoreAssets`, and `_filePath`
functions to produce cleaner godocs.

## License Info
This software is based on the NOAA World Magnetic Model.
The source code in this project is not based on the source code provided by NOAA, but on the
equations provided in the World Magnetic Model reference paper.

The WMM source code is not subject to copyright protection: https://www.ngdc.noaa.gov/geomag/WMM/license.shtml

The WMM source code is in the public domain and not licensed or under copyright. The information and software may be used freely by the public. As required by 17 U.S.C. 403, third parties producing copyrighted works consisting predominantly of the material produced by U.S. government agencies must provide notice with such work(s) identifying the U.S. Government material incorporated and stating that such material is not subject to copyright protection.
