package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
)

const (
	dbHost     = "mssql-service.default.svc.cluster.local"
	dbPort     = 1433
	dbUser     = "sa"
	dbPassword = "YourStrong!Passw0rd" // Change this to your password
	dbName     = "LoanCRM"
	dbSchema   = "dbo"
)

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		dbHost, dbUser, dbPassword, dbPort, dbName)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error connecting to database:", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot ping database:", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	// Define record counts based on 1,000,000 opportunities
	numOpportunities := 1000000
	numBusinessPartners := int(float64(numOpportunities) * 1.25) // 1.25 partners per opportunity on average
	numCases := int(float64(numOpportunities) * 0.4)             // ~40% of opportunities have a related case

	// Generate and insert data
	log.Println("Generating and inserting Business Partners...")
	businessPartners := generateAndInsertBusinessPartners(db, numBusinessPartners)
	log.Println("Done. Total Business Partners:", len(businessPartners))

	log.Println("Generating and inserting Opportunities...")
	opportunities := generateAndInsertOpportunities(db, numOpportunities, businessPartners)
	log.Println("Done. Total Opportunities:", len(opportunities))

	log.Println("Generating and inserting Opportunity-Partner relationships...")
	generateAndInsertOpportunityPartners(db, opportunities, businessPartners)
	log.Println("Done.")

	log.Println("Generating and inserting Cases...")
	generateAndInsertCases(db, numCases, businessPartners, opportunities)
	log.Println("Done. Total Cases:", numCases)

	log.Println("Database population complete.")
}

func generateAndInsertBusinessPartners(db *sql.DB, num int) []uuid.UUID {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback() // Rollback if not committed

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s.tbl_BusinessPartners (PartnerGUID, FirstName, LastName, Email, PhoneNumber) VALUES (?, ?, ?, ?, ?)", dbSchema))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	partnerGUIDs := make([]uuid.UUID, num)
	firstNames := []string{"John", "Jane", "Robert", "Susan", "Michael", "Linda"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis"}

	for i := 0; i < num; i++ {
		guid := uuid.New()
		partnerGUIDs[i] = guid
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@example.com", firstName, lastName, i)
		phoneNumber := fmt.Sprintf("555-%d", 10000000+i)

		_, err = stmt.Exec(guid, firstName, lastName, email, phoneNumber)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
	return partnerGUIDs
}

func generateAndInsertOpportunities(db *sql.DB, num int, partners []uuid.UUID) []uuid.UUID {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s.tbl_Opportunities (OpportunityGUID, ProcessType, Status, RequestedAmount, CreatedDate, ClosingDate) VALUES (?, ?, ?, ?, ?, ?)", dbSchema))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	opportunityGUIDs := make([]uuid.UUID, num)
	processTypes := []string{"Mortgage", "Auto Loan", "Personal Loan"}
	statuses := []string{"In Progress", "Approved", "Rejected"}

	for i := 0; i < num; i++ {
		guid := uuid.New()
		opportunityGUIDs[i] = guid
		processType := processTypes[rand.Intn(len(processTypes))]
		status := statuses[rand.Intn(len(statuses))]
		amount := float64(rand.Intn(1000000-1000) + 1000)
		createdDate := time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour)
		closingDate := createdDate.Add(time.Duration(rand.Intn(60)) * 24 * time.Hour)

		_, err = stmt.Exec(guid, processType, status, amount, createdDate, closingDate)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
	return opportunityGUIDs
}

func generateAndInsertOpportunityPartners(db *sql.DB, opportunities, partners []uuid.UUID) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s.tbl_OpportunityPartners (OpportunityGUID, PartnerGUID, PartnerFunction, IsPrimary) VALUES (?, ?, ?, ?)", dbSchema))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, oppGUID := range opportunities {
		// Each opportunity has at least one primary partner
		primaryPartnerGUID := partners[rand.Intn(len(partners))]
		_, err := stmt.Exec(oppGUID, primaryPartnerGUID, "Main Borrower", 1)
		if err != nil {
			log.Fatal(err)
		}

		// Add a co-borrower for some opportunities (~40% chance)
		if rand.Float64() < 0.4 {
			coBorrowerGUID := partners[rand.Intn(len(partners))]
			// Ensure co-borrower is not the same as the main borrower
			for coBorrowerGUID == primaryPartnerGUID {
				coBorrowerGUID = partners[rand.Intn(len(partners))]
			}
			_, err := stmt.Exec(oppGUID, coBorrowerGUID, "Co-Borrower", 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	tx.Commit()
}

func generateAndInsertCases(db *sql.DB, num int, partners, opportunities []uuid.UUID) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s.tbl_Cases (CaseGUID, CaseType, Status, Summary, AssignedTo, CreatedDate, ClosedDate, PartnerGUID, OpportunityGUID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", dbSchema))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	caseTypes := []string{"Payment Inquiry", "Complaint", "Document Request"}
	statuses := []string{"Open", "Resolved"}
	assignedTos := []string{"Susan Miller", "John Davis", "Lisa White"}

	for i := 0; i < num; i++ {
		guid := uuid.New()
		caseType := caseTypes[rand.Intn(len(caseTypes))]
		status := statuses[rand.Intn(len(statuses))]
		summary := fmt.Sprintf("Generated summary for %s.", caseType)
		assignedTo := assignedTos[rand.Intn(len(assignedTos))]
		createdDate := time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour)
		var closedDate sql.NullTime
		if status == "Resolved" {
			closedDate.Time = createdDate.Add(time.Duration(rand.Intn(15)) * 24 * time.Hour)
			closedDate.Valid = true
		}

		// Link case to an existing partner
		partnerGUID := partners[rand.Intn(len(partners))]

		// Link case to an existing opportunity (~60% chance)
		var opportunityGUID sql.Null[uuid.UUID]
		if rand.Float64() < 0.6 {
			opportunityGUID.V = opportunities[rand.Intn(len(opportunities))]
			opportunityGUID.Valid = true
		}

		_, err = stmt.Exec(guid, caseType, status, summary, assignedTo, createdDate, closedDate, partnerGUID, opportunityGUID)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}
