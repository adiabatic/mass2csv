# mass2csv

`mass2csv` reads your exported HealthKit data and outputs your dated weight measurements as comma-separated values, in pounds.

## Usage

```fish
mass2csv filename
```

`filename` can either be the `export.xml` or the entire `export.zip`. Note that it won’t know what to do with files formatted like `export_cda.xml`.

`mass2csv` outputs comma-separated values to standard output, so you’ll probably want to redirect that to a file. To do so, run `mass2csv` like this in the shell of your choice:

```fish
mass2csv export.zip > weights.csv
```

Then open `weights.csv` in your favorite spreadsheet software.

## Exit status

0 on success, 1 on failure. Failures are accompanied by a message printed to the standard error stream.

## Notes

- Entries aren’t guaranteed to be sorted chronologically. You’ll want to sort them in your favorite spreadsheet software.
- I’ve never entered masses in kilograms into HealthKit, so I have no idea if that particular code path is correct.
- `mass2csv` isn’t equipped to read `export_cda.xml` files. If you try to pass one of those to `mass2csv`, it will output only the CSV header.

## License

Public domain.
