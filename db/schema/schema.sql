--     Database:  wreis (Wertz Real Estate Investment Services)
--
--     Field names are camel case
--     Money values are all stored as DECIMAL(19,4)

DROP DATABASE IF EXISTS wreis;
CREATE DATABASE wreis;
USE wreis;
GRANT ALL PRIVILEGES ON wreis.* TO 'ec2-user'@'localhost';
set GLOBAL sql_mode='ALLOW_INVALID_DATES';


CREATE TABLE Property (
    PRID BIGINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(256) NOT NULL DEFAULT '',
    YearsInBusiness SMALLINT NOT NULL DEFAULT 0,
    ParentCompany VARCHAR(256) NOT NULL DEFAULT '',
    URL VARCHAR(1028) NOT NULL DEFAULT '',                  -- web address
    Symbol VARCHAR(128) NOT NULL DEFAULT '',                -- Stock Symbol and board
    Price DECIMAL(19,4) NOT NULL DEFAULT 0,
    DownPayment DECIMAL(19,4) NOT NULL DEFAULT 0,
    RentableArea BIGINT NOT NULL DEFAULT 0,
    RentableAreaUnits SMALLINT NOT NULL DEFAULT 0,          -- 0 = sqft, 1 = acres,
    LotSize BIGINT NOT NULL DEFAULT 0,
    LotSizeUnits SMALLINT NOT NULL DEFAULT 0,               -- 0 = sqft, 1 = acres,
    CapRate FLOAT NOT NULL DEFAULT 0,                       -- percentage
    AvgCap FLOAT NOT NULL DEFAULT 0,                        -- percentage
    BuildDate DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00', -- Date the property was built, if applicable
    FLAGS BIGINT NOT NULL DEFAULT 0,                        /* 1<<0  Drive Through?  0 = no, 1 = yes
                                                               1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
                                                               1<<2  Right Of First Refusal: 0 = no, 1 = yes
                                                            */
    Ownership SMALLINT NOT NULL DEFAULT 0,                  -- 0 = fee simple, 1 = leasehold
    TenantTradeName VARCHAR(256) NOT NULL DEFAULT '',       -- trade name of business
    LeaseGuarantor SMALLINT NOT NULL DEFAULT 0,             -- 0 = corporate, 1 = franchise, 2 = individual
    LeaseType SMALLINT NOT NULL DEFAULT 0,                  -- 0 = Absolute NNN, 1 = Double Net, 2 = Triple Net, 3 = Gross
    DeliveryDt DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',  -- GMT datetime
    OriginalLeaseTerm BIGINT NOT NULL DEFAULT 0,            -- Duration
    RentCommencementDt DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',    -- GMT datetime
    LeaseExpirationDt DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',      -- GMT datetime
    TermRemainingOnLease BIGINT NOT NULL DEFAULT 0,         -- Duration
    TermRemainingOnLeaseUnits SMALLINT NOT NULL DEFAULT 0,  -- 0 = months, 1 = Years
    ROLID BIGINT NOT NULL DEFAULT 0,                        -- ID of associated Renew Options
    RSLID BIGINT NOT NULL DEFAULT 0,                        -- ID of associated Rent Steps

    Address VARCHAR(100) NOT NULL DEFAULT '',               -- property address
    Address2 VARCHAR(100) NOT NULL DEFAULT '',
    City VARCHAR(100) NOT NULL DEFAULT '',
    State CHAR(25) NOT NULL DEFAULT '',
    PostalCode VARCHAR(100) NOT NULL DEFAULT '',
    Country VARCHAR(100) NOT NULL DEFAULT '',

    LLResponsibilities VARCHAR(2048) NOT NULL DEFAULT '',   -- Is this enough characters
    NOI DECIMAL(19,4) NOT NULL DEFAULT 0,                   -- Net Operating Income

    HQAddress VARCHAR(100) NOT NULL DEFAULT '',             -- address of headquarters only City/State are required
    HQAddress2 VARCHAR(100) NOT NULL DEFAULT '',
    HQCity VARCHAR(100) NOT NULL DEFAULT '',
    HQState CHAR(25) NOT NULL DEFAULT '',
    HQPostalCode VARCHAR(100) NOT NULL DEFAULT '',
    HQCountry VARCHAR(100) NOT NULL DEFAULT '',

    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTS TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (PRID)
);

CREATE TABLE RenewOptions (
    ROLID BIGINT NOT NULL AUTO_INCREMENT,                   -- Renew Options List ID
    /*
    ++  bit 0 - There are 2 fundamental ways in which Renew Options are specified.
    **          bit 0 set to 0 means that each RenewOption record specifies an absolute
    **          date.  Bit 0 set to 1 means that each RenewOption record specifies
    **          a count of years past commencement.  Examples:
    **
    **    ----------- bit 0 = 0 -----------     ----------- bit 0 = 1 -----------
    **     String
    **     Option                 Annual        Option       Option     Annual
    **     Period                 Rent          Year         Period     Rent
    **    ---------------------------------     ------------------------------------
    **     Year 1-4               183,568.85    7/4/2024     1          109,709.45
    **     Year 5-9               201,925.74    7/4/2025     1          111,903.63
    **     Year 10                222,118.32    7/4/2026     2          116.424.54
    **                                          ...
    */
    FLAGS BIGINT NOT NULL DEFAULT 0,                        -- 1<<0 = 0 = counts, 1 = dates (see comment above)
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTS TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (ROLID)
);

CREATE TABLE RenewOption (
    ROID BIGINT NOT NULL AUTO_INCREMENT,                    -- A Renew Option, part of a list
    ROLID BIGINT NOT NULL DEFAULT 0,                        -- Renew Options List ID to which this RO belongs
    Dt DATE NOT NULL DEFAULT '1970-01-01 00:00:00',         -- Date that the rent went into effect, valid only when ROLID FLAGS bit 0 = 1
    Opt VARCHAR(128) NOT NULL DEFAULT '',                   -- option period comment
    Rent DECIMAL(19,4) NOT NULL DEFAULT 0,                  -- Monthly Rent Amount
    FLAGS BIGINT NOT NULL DEFAULT 0,                        -- 1<<0 = 0 = counts, 1 = dates (see comment above)
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTS TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (ROID)
);

CREATE TABLE RentSteps (
    RSLID BIGINT NOT NULL AUTO_INCREMENT,                   -- RentStep List ID
    FLAGS BIGINT NOT NULL DEFAULT 0,                        -- 1<<0 = 0 = count, 1 = dates -- See comment for RenewOptions FLAGS
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTS TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (RSLID)
);

CREATE TABLE RentStep (
    RSID BIGINT NOT NULL AUTO_INCREMENT,                    -- A Rent Step, part of a list
    RSLID BIGINT NOT NULL DEFAULT 0,                        -- RentStep List ID to which this RS belongs
    Dt DATE NOT NULL DEFAULT '1970-01-01 00:00:00',         -- Date that the rent went into effect, valid only when ROLID FLAGS bit 0 = 1
    Opt VARCHAR(128) NOT NULL DEFAULT '',                   -- option period comment
    Rent DECIMAL(19,4) NOT NULL DEFAULT 0,                  -- Rent commencement date
    FLAGS BIGINT NOT NULL DEFAULT 0,                        -- 1<<0 = 0 = options, 1 = dates (see comment above)
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTS TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (RSID)
);
