
CREATE TABLE IF NOT EXISTS Users
(
    IdUser      INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    Progression INTEGER                           NOT NULL,
    Name        TEXT(30)                          NOT NULL,
    Email       TEXT(255)                         NOT NULL,
    Pwd         TEXT(255)                         NOT NULL,
    Score       INTEGER                           NOT NULL,
    FOREIGN KEY (Progression) REFERENCES ExerciceDone (Progression)
);

CREATE TABLE IF NOT EXISTS Exercices
(
    IdExo       INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    Progression INTEGER                           NOT NULL,
    NameExo     TEXT(30)                          NOT NULL,
    Description TEXT(255)                         NOT NULL,
    Difficulte  INTEGER                           NOT NULL
);

CREATE TABLE IF NOT EXISTS ExerciceDone
(
    IdExo       INTEGER  NOT NULL,
    Progression INTEGER  NOT NULL,
    FOREIGN KEY (Progression) REFERENCES Exercices (Progression),
    FOREIGN KEY (IdExo) REFERENCES Exercices (IdExo)
);
