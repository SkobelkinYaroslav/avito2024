CREATE TABLE Users (
                       ID VARCHAR(255) PRIMARY KEY,
                       Email VARCHAR(255) NOT NULL UNIQUE,
                       Password VARCHAR(255) NOT NULL,
                       User_Type VARCHAR(50) NOT NULL
);

CREATE TABLE Token(
                ID SERIAL PRIMARY KEY,
                User_ID VARCHAR(255) NOT NULL,
                Token VARCHAR(255) NOT NULL,
                FOREIGN KEY (User_ID) REFERENCES Users(ID)

);

CREATE TABLE Houses (
                        ID SERIAL PRIMARY KEY,
                        Address VARCHAR(255) NOT NULL,
                        Year INT NOT NULL,
                        Developer VARCHAR(255),
                        CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Flats (
                       ID SERIAL PRIMARY KEY,
                       House_ID INT NOT NULL,
                       Price INT NOT NULL,
                       Rooms INT NOT NULL,
                       Status VARCHAR(50) NOT NULL,
                       FOREIGN KEY (House_ID) REFERENCES Houses(ID)
);