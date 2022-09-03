package dbengine

import "errors"

// InsertFile creates a new record on file database
// from given File struct
func InsertFile(file File) error {
	if err := ensureConnection(); err != nil {
		return err
	}

	statement := `INSERT INTO files(
		filename,
		password,
		upload_date,
		file_size,
		viewcount,
		owner_id
	) VALUES(?, ?, ?, ?, ?, ?);`

	if _, err := DbConnection.Exec(statement, file.Filename, file.Password,
		file.UploadDate, file.FileSize,
		file.ViewCount, file.OwnerId); err != nil {
		return err
	}

	return nil
}

// GetPasswordHash returns password field from database
func GetPasswordHash(filename string) (pwd string, err error) {
	if err = ensureConnection(); err != nil {
		return "", err
	}

	row := DbConnection.QueryRow("SELECT password FROM files WHERE filename = ?;", filename)
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
  - FileSize
  - ViewCount
*/
func GetFile(filename string) (f File, err error) {
	if err = ensureConnection(); err != nil {
		return File{}, err
	}

	row := DbConnection.QueryRow("SELECT filename, upload_date, file_size, viewcount FROM files WHERE filename = ?",
		filename)
	if err = row.Scan(&f.Filename, &f.UploadDate,
		&f.FileSize, &f.ViewCount); err != nil {
		return File{}, err
	}

	return f, nil
}

// IncrementFileViewCount increments file's viewcount
// field by 1 in database
func IncrementFileViewCount(filename string) error {
	if err := ensureConnection(); err != nil {
		return err
	}

	statement := "UPDATE files SET viewcount = viewcount + 1 WHERE filename = ?"
	if _, err := DbConnection.Exec(statement, filename); err != nil {
		return err
	}

	return nil
}

func ensureConnection() error {
	if DbConnection == nil {
		return errors.New("database is not connected")
	}
	return nil
}
