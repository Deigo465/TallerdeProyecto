-- Adding test data for profiles
INSERT INTO profiles (first_name, mother_last_name, father_last_name, document_number, gender, phone, contact_email, date_of_birth, cmp, specialty, role) VALUES
('Staff', 'Doe', 'Smith', '12345678', 'Male', '123456789', "staff@example.com", '1970-01-15', 'CMP001', 'Cardiology', 'STAFF');
-- ('Doctor', 'Jones', 'White', '12348765', 'Female', '345678123',  "doctor@example.com", '1990-03-05', '', '', 'DOCTOR'),
-- -- Add patients
-- ('Patient', 'Jones', 'White', '12348765', 'Female', '345678123', "contact@email.com", '1990-03-05', '', '', 'PATIENT'),
-- ('Michael', 'Brown', 'Clark', '56781234', 'Male', '456789234', "contact@email.com", '1985-04-10', '', '', 'PATIENT'),
-- ('Emma', 'Wilson', 'Martin', '43561278', 'Female', '567891345', "contact@email.com", '1995-05-15', '', '', 'PATIENT'),
-- ('Oliver', 'Taylor', 'Lee', '87651234', 'Male', '678912345', "contact@email.com", '1982-06-20', '', '', 'PATIENT'),
-- ('Isabella', 'Moore', 'Walker', '76543218', 'Female', '912345678', "contact@email.com", '1987-09-05', '', '', 'PATIENT');

-- Adding test data for health_centers
INSERT INTO health_centers (name, district, address) VALUES
('Health Center 1', 'District 1', '123 Main St'),
('Health Center 2', 'District 2', '456 Broadway'),
('Health Center 3', 'District 3', '789 High St');

-- Adding test data for users
-- All belong to the health center 1
INSERT INTO users (email, password, health_center_id, profile_id) VALUES
('staff@example.com', '123456', 1, 1)
-- ('doctor@example.com', '123456', 1, 2);

-- Adding test data for records
-- INSERT INTO records (body, created_at, updated_at, patient_record_id, doctor_record_id, specialty) VALUES
-- ('Record 1 body text', '2023-01-01 10:00:00', '2023-01-01 10:00:00', 3, 1, 'Odontologia'),


-- Adding test data for files
-- INSERT INTO files (url, name, filesize, mimetype, record_id) VALUES
-- ('http://example.com/file1', 'file1', "1000", 'text/plain', 1)

-- Adding test data for appointments
-- INSERT INTO appointments (specialty, status, date_appointment, doctor_id, patient_id, description) VALUES
-- ('Dermatology', 1, '2024-06-01', 1,6,'estoy mal'),
-- ('Dermatology', 1, '2024-06-02', 1,6,'estoy mal'),
-- ('Gynecology', 1, '2024-06-03', 1,6,'estoy mal'),
-- ('Pediatrics', 2, '2024-06-04', 1,6,'estoy mal'),
-- ('Cardiology', 1, '2024-05-05', 1,6,'estoy mal'),
-- ('Dermatology', 2, '2024-06-06', 1,6,'estoy mal'),
-- ('Gynecology', 1, '2024-06-07', 1,6,'estoy mal'),
-- ('Pediatrics', 2, '2024-06-08', 1,6,'estoy mal'),
-- ('Cardiology', 1, '2024-06-09', 1,6,'estoy mal'),
-- ('Dermatology', 2, '2024-06-10', 1,6,'estoy mal');


-- INSERT INTO sessions (user_id, token, created_at, updated_at) VALUES
-- (1, 'STAFF', '2025-01-01 10:00:00', '2023-01-01 10:00:00'),
-- (2, 'DOCTOR', '2025-01-01 10:00:00', '2023-01-01 10:00:00'),
-- (3, 'PATIENT', '2025-01-01 10:00:00', '2023-01-01 10:00:00');