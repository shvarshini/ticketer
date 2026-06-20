ALTER TABLE seats ADD CONSTRAINT seats_screen_id_row_number_key UNIQUE (screen_id, row, number);
