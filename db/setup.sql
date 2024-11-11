CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    first_name VARCHAR(50),
    mother_last_name VARCHAR(50),
    father_last_name VARCHAR(50),
    document_number VARCHAR(8),
    gender VARCHAR(10),
    phone VARCHAR(9),
    contact_email VARCHAR(50),
    date_of_birth DATE,
    cmp VARCHAR(50),
    specialty VARCHAR(50),
    role VARCHAR(10) CHECK(role IN ('STAFF', 'DOCTOR', 'PATIENT', 'ADMIN'))
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    email VARCHAR(50),
    password VARCHAR(50),
    health_center_id INTEGER NOT NULL,
    profile_id INTEGER NOT NULL,
    FOREIGN KEY (health_center_id) REFERENCES health_centers(id),
    FOREIGN KEY (profile_id) REFERENCES profiles(id)
);


CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    token VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);




CREATE TABLE IF NOT EXISTS health_centers (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    name VARCHAR(50),
    district VARCHAR(50),
    address VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS records (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    body TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    patient_record_id INTEGER NOT NULL,
    doctor_record_id INTEGER NOT NULL,
    specialty VARCHAR(50) NOT NULL,
    FOREIGN KEY (patient_record_id) REFERENCES profiles(id),
    FOREIGN KEY (doctor_record_id) REFERENCES profiles(id)
);

CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    url VARCHAR(255),
    name VARCHAR(50),
    filesize VARCHAR(10),
    mimetype VARCHAR(10),
    record_id INTEGER NOT NULL,
    FOREIGN KEY (record_id) REFERENCES records(id)
);

CREATE TABLE IF NOT EXISTS appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    specialty VARCHAR(100),
    status INTEGER CHECK(status IN (0, 1, 2, 3, 4)), 
    date_appointment DATE,
    doctor_id INTEGER NOT NULL,
    patient_id INTEGER NOT NULL,
    description VARCHAR(255) NOT NULL,
    FOREIGN KEY (doctor_id) REFERENCES profiles(id)
    FOREIGN KEY (patient_id) REFERENCES profiles(id)

);
