### MSSQL Database Schematic

```tsql
CREATE TABLE dbo.Messages
(
    Id        CHAR(36)                      NOT NULL
        CONSTRAINT PK_Messages_Id
            PRIMARY KEY,
    Recipient VARCHAR(25)                   NOT NULL,
    Content   VARCHAR(255)                  NOT NULL,
    Status    VARCHAR(25) DEFAULT 'Unsent'  NOT NULL
        CONSTRAINT chk_status
            CHECK ([Status] = 'Unsent' OR [Status] = 'Sent'),
    CreatedAt DATETIME    DEFAULT GETDATE() NOT NULL,
    SentAt    DATETIME
);
GO
