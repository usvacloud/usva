package dbengine

// InsertFile creates a new record on file database
// from given File struct
func InsertFile(file File) error {
	if _, err := DbConnection.Exec(insertFileQuery, file.Filename, file.Password,
		file.UploadDate, file.IsEncrypted,
		file.ViewCount, file.OwnerId); err != nil {
		return err
	}

	return nil
}

// DeleteFile removes file metadata from database
func DeleteFile(filename string) error {
	if _, err := DbConnection.Exec(deleteFileQuery, filename); err != nil {
		return err
	}

	return nil
}

// GetPasswordHash returns password field from database
func GetPasswordHash(filename string) (pwd string, err error) {
	row := DbConnection.QueryRow(getPasswordQuery, filename)
	err = row.Scan(&pwd)

	return pwd, err
}

/*
GetFile returns File struct populated with file's
metadata from database on successfull request

Note: only following fields are included, thus any other
fields will remain empty:
  - Filename
  - UploadDate
  - IsEncrypted
  - ViewCount
*/
func GetFile(filename string) (f File, err error) {
	row := DbConnection.QueryRow(getFileQuery, filename)
	if err = row.Scan(&f.Filename, &f.UploadDate, &f.IsEncrypted,
		&f.ViewCount); err != nil {
		return File{}, err
	}

	return f, nil
}

// IncrementFileViewCount increments file's viewcount
// field by 1 in database
func IncrementFileViewCount(filename string) error {
	if _, err := DbConnection.Exec(incrementFileViewCountQuery, filename); err != nil {
		return err
	}

	return nil
}
