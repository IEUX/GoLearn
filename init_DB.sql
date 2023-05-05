CREATE TABLE IF NOT EXISTS User(
    ID_User INTEGER PRIMARY KEY,
    Progression INTEGER ,
    Username varchar(30),
    Email varchar(255),
    Password varchar(255),
    Exp INTEGER,
    FOREIGN KEY (Progression) REFERENCES Exercise(ID_Exercise)
);

CREATE TABLE IF NOT EXISTS Exercise(
    ID_Exercise INTEGER PRIMARY KEY,
    Title varchar(255),
    Prompt TEXT,
    Difficulty INTEGER
);

CREATE TABLE IF NOT EXISTS Solution(
    ID_Solution INTEGER PRIMARY KEY,
    ID_Exercise INTEGER ,
    Solution TEXT,
    FOREIGN KEY (ID_Exercise) REFERENCES Exercise(ID_Exercise)
);