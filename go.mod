module github.com/kristianvalind/go-netdicom

go 1.16

require (
	github.com/stretchr/testify v1.7.0
	github.com/suyashkumar/dicom v1.0.3
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/tools v0.1.3 // indirect
)

replace (
	github.com/suyashkumar/dicom v1.0.3 => ../dicom
)
