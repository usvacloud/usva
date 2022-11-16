package dbengine

// InsertFile creates a new record on file database
// from given File struct
func InsertFile(file File) error {
	_, err := DbConnection.Exec(insertFileQuery, file.FileUUID, file.Title, file.Uploader,
		file.PasswordHash, file.UploadDate, file.IsEncrypted)
	return err
}

// DeleteFile removes file metadata from database
func DeleteFile(filename string) error {
	_, err := DbConnection.Exec(deleteFileQuery, filename)
	return err
}

// TryDeleteFile is DeleteFile, but ignores any errors.
func TryDeleteFile(filename string) {
	_ = DeleteFile(filename)
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
  - Title
  - UploadDate
  - IsEncrypted
  - ViewCount
*/
func GetFile(filename string) (f File, err error) {
	row := DbConnection.QueryRow(getFileInformationQuery, filename)
	err = row.Scan(&f.FileUUID, &f.Title, &f.UploadDate, &f.IsEncrypted, &f.ViewCount)
	return f, err
}

// IncrementFileViewCount increments file's viewcount
// field by 1 in database
func IncrementFileViewCount(filename string) error {
	_, err := DbConnection.Exec(incrementFileViewCountQuery, filename)
	return err
}
