CREATE DATABASE [LoanCRM];
GO

USE LoanCRM;
GO

-- Create tbl_Opportunities
CREATE TABLE tbl_Opportunities (
    OpportunityGUID UNIQUEIDENTIFIER PRIMARY KEY,
    ProcessType VARCHAR(50) NOT NULL,
    Status VARCHAR(50) NOT NULL,
    Description VARCHAR(255),
    RequestedAmount DECIMAL(18, 2),
    CreatedDate DATETIME,
    ClosingDate DATETIME
);
GO

-- Create tbl_BusinessPartners
CREATE TABLE tbl_BusinessPartners (
    PartnerGUID UNIQUEIDENTIFIER PRIMARY KEY,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    Email VARCHAR(150),
    PhoneNumber VARCHAR(20)
);
GO

-- Create tbl_OpportunityPartners (Junction Table)
CREATE TABLE tbl_OpportunityPartners (
    OpportunityGUID UNIQUEIDENTIFIER,
    PartnerGUID UNIQUEIDENTIFIER,
    PartnerFunction VARCHAR(50) NOT NULL,
    IsPrimary BIT NOT NULL,
    CONSTRAINT PK_OpportunityPartners PRIMARY KEY (OpportunityGUID, PartnerGUID),
    CONSTRAINT FK_OpportunityPartners_Opportunities FOREIGN KEY (OpportunityGUID) REFERENCES tbl_Opportunities(OpportunityGUID),
    CONSTRAINT FK_OpportunityPartners_Partners FOREIGN KEY (PartnerGUID) REFERENCES tbl_BusinessPartners(PartnerGUID)
);
GO

-- Create tbl_Cases
CREATE TABLE tbl_Cases (
    CaseGUID UNIQUEIDENTIFIER PRIMARY KEY,
    CaseType VARCHAR(50) NOT NULL,
    Status VARCHAR(50) NOT NULL,
    Summary VARCHAR(255),
    AssignedTo VARCHAR(100),
    CreatedDate DATETIME,
    ClosedDate DATETIME,
    PartnerGUID UNIQUEIDENTIFIER,
    OpportunityGUID UNIQUEIDENTIFIER,
    CONSTRAINT FK_Cases_Partner FOREIGN KEY (PartnerGUID) REFERENCES tbl_BusinessPartners(PartnerGUID),
    CONSTRAINT FK_Cases_Opportunity FOREIGN KEY (OpportunityGUID) REFERENCES tbl_Opportunities(OpportunityGUID)
);
GO