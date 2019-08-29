

# Usage

    --admin-db string
        db port (default "postgres")
        
    --admin-pass-file string
        file with admin passwords, one per line
        (not required if DB_ADMIN_PASSWORDS env var is set accordingly)
        
    --admin-user string
        db port (default "postgres")
        
    --engine string
        db engine (default "postgresql")
        NOTE: this is a placeholder option, only the default
              'postgresql' option is supported for now
        
    --host string
        db host (default "localhost")
        
    --port string
        db port (default "5432")
        
    --user string
        db user to create (if not exists)
        
    --user-db string
        user db to create (if not exists)
        
    --user-pass-file string
        file with user passwords, one per line
        (not required if DB_USER_PASSWORDS env var is set accordingly)


## Example

NOTE: at least one line of `--admin-pass-file` in this example contains the actual admin user password
(and real cases must adhere to this requirement)

`db-maker.go --admin-pass-file adminpass.txt --user user1 --user-db db1 --user-pass-file userpass.txt`

this command fetches passwords from `adminpass.txt` and `userpass.txt`, sets `--admin-user` 
password to the first non-empty line of the `adminpass.txt` file, creates `--user`
with the first non-empty line of `--user-pass-file` as the user's password, creates `--user-db`
with owner `--user`, updates `--user` password to the first non empty line of `--user-pass-file`
if it is not his password yet

`--admin-pass-file` and `--user-pass-file` may be substituted with:
DB_ADMIN_PASSWORDS=`printf "PASS1\nPASS2\n..."`
DB_USER_PASSWORDS=...
