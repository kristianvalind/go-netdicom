module github.com/kristianvalind/go-netdicom

go 1.16

require (
	github.com/stretchr/testify v1.7.0
	github.com/suyashkumar/dicom v1.0.3
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/suyashkumar/dicom v1.0.3 => ../dicom
