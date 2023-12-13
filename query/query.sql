-- DDL untuk tabel Heroes
CREATE TABLE Heroes (
    ID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(255) NOT NULL,
    Universe VARCHAR(255) NOT NULL,
    Skill VARCHAR(255),
    ImageURL VARCHAR(255)
);

-- DDL untuk tabel Villain
CREATE TABLE Villain (
    ID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(255) NOT NULL,
    Universe VARCHAR(255) NOT NULL,
    ImageURL VARCHAR(255)
);

-- DDL untuk tabel CrimeEvent
CREATE TABLE CrimeEvent (
    ID INT PRIMARY KEY AUTO_INCREMENT,
    HeroID INT,
    VillainID INT,
    Description TEXT,
    DateTime DATETIME,
    FOREIGN KEY (HeroID) REFERENCES Heroes(ID),
    FOREIGN KEY (VillainID) REFERENCES Villain(ID)
);

CREATE TABLE item (
    ID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(255) NOT NULL,
    ItemCode VARCHAR(50) NOT NULL,
    Stock INT NOT NULL,
    Description VARCHAR(255),
    Status VARCHAR(50) NOT NULL CHECK (Status IN ('Active', 'Broken'))
);

-- DML untuk insert data ke tabel Heroes
INSERT INTO Heroes (ID, Name, Universe, Skill, ImageURL)
VALUES
    (1, 'Superman', 'DC', 'Super strength, flight', 'superman.jpg'),
    (2, 'Spider-Man', 'Marvel', 'Wall-crawling, web-shooting', 'spiderman.jpg'),
    (3, 'Wonder Woman', 'DC', 'Super strength, Lasso of Truth', 'wonderwoman.jpg'),
    (4, 'Iron Man', 'Marvel', 'Powered armor suit', 'ironman.jpg'),
    (5, 'Black Widow', 'Marvel', 'Espionage, martial arts', 'blackwidow.jpg');

-- DML untuk insert data ke tabel Villain
INSERT INTO Villain (ID, Name, Universe, ImageURL)
VALUES
    (1, 'Lex Luthor', 'DC', 'lexluthor.jpg'),
    (2, 'Green Goblin', 'Marvel', 'greengoblin.jpg'),
    (3, 'Cheetah', 'DC', 'cheetah.jpg'),
    (4, 'Thanos', 'Marvel', 'thanos.jpg'),
    (5, 'Red Skull', 'Marvel', 'redskull.jpg');

-- DML untuk insert data ke tabel CrimeEvent
INSERT INTO CrimeEvent (ID, HeroID, VillainID, Description, DateTime)
VALUES
    (1, 1, 1, 'Lex Luthor robs a bank', '2023-12-13 08:30:00'),
    (2, 2, 2, 'Green Goblin attacks Times Square', '2023-12-14 15:45:00'),
    (3, 3, 3, 'Cheetah kidnaps a diplomat', '2023-12-15 12:00:00'),
    (4, 4, 4, 'Thanos threatens the world', '2023-12-16 18:20:00'),
    (5, 5, 5, 'Red Skull plots world domination', '2023-12-17 09:10:00');

INSERT INTO item (Name, ItemCode, Stock, Description, Status) VALUES
    ('Item 1', 'CODE001', 50, 'Description 1', 'Active'),
    ('Item 2', 'CODE002', 30, 'Description 2', 'Broken'),
    ('Item 3', 'CODE003', 20, 'Description 3', 'Active'),
    ('Item 4', 'CODE004', 40, 'Description 4', 'Broken'),
    ('Item 5', 'CODE005', 10, 'Description 5', 'Active'),
    ('Item 6', 'CODE006', 25, 'Description 6', 'Broken'),
    ('Item 7', 'CODE007', 35, 'Description 7', 'Active'),
    ('Item 8', 'CODE008', 15, 'Description 8', 'Broken'),
    ('Item 9', 'CODE009', 45, 'Description 9', 'Active'),
    ('Item 10', 'CODE010', 5, 'Description 10', 'Broken');

