INSERT INTO Houses (Address, Year, Developer)
VALUES ('123 Main St', 2000, 'Real Estate Developer');

INSERT INTO Flats (House_ID, Price, Rooms, Status)
VALUES
    (1, 300000, 3, 'created'),
    (1, 250000, 2, 'approved'),
    (1, 400000, 4, 'on moderation'),
    (1, 150000, 1, 'declined'),
    (1, 500000, 5, 'approved'),
    (1, 350000, 3, 'created');
