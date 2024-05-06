create table Exercise
(
    ID_Exercise        INTEGER
        primary key,
    Title              varchar(255),
    Prompt             TEXT,
    Difficulty         INTEGER,
    function_structure TEXT
);

create table Solution
(
    ID_Solution INTEGER
        primary key,
    ID_Exercise INTEGER
        references Exercise,
    Solution    TEXT,
    Test        TEXT
);

create table User
(
    ID_User     INTEGER
        primary key,
    Progression INTEGER
        references Exercise,
    Username    varchar(30),
    Email       varchar(255),
    Password    varchar(255),
    Exp         INTEGER
);
