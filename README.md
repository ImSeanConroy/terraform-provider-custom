# Terraform Provider Custom

Custom Terraform provider designed to interact with a custom API.

## Run Tests
```
cd internal/provider && TF_ACC=1 go test -v .
```

## Generate Documentation
```
cd tools && go generate ./...
```

## Future Improvements

[] Add Group Association
[] Add Automated Resource and Datasource Genertion Development Process