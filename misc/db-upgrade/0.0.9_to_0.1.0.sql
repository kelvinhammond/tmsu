-- idx_file_path is redundant and an overhead as the unique constraint creates an identical index
DROP INDEX IF EXISTS idx_file_path;

-- new size column on the file table
ALTER TABLE file ADD COLUMN size INTEGER NOT NULL DEFAULT 0;

-- invalid dates cause a problem for go-sqlite 3 so need to set them to epoch.
UPDATE file SET mod_time = '1970-01-01 00:00:00' WHERE mod_time = '0000-00-00 00:00:00';

-- the new size field does not get repaired by the 'repair' command unless it thinks the file has changed so
-- set the mod_time to epoch for all files with a zero size to force the 'repair' command to update them.
UPDATE file SET mod_time = '1970-01-01 00:00:00' WHERE size = 0;
