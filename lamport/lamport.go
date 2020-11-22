package lamport

import (
	"common"
	"config"
	"net"
	"sync/atomic"
)

/* VARIABLES */

// Nombre de processus connectés
var connectedProcess int32 = 0
// Nombre de processus total
var nbProcess int32 = 0

// Identifiant du processus actuel
var idCurrentProcess int

// Estampille du processus actuel
type Clock struct {
	Timestamp uint64
}
var clock = Clock{Timestamp: 0}
// Si le processus actuel à le droit à la SC
var canSc = false

// Stock la dernière requête émise par chaque processus
var processRequest = make(map[int]ConnectionProcess)
// Structure représentant un processus avec sa référence, sa connexion
//et sa dernière requête
type ConnectionProcess struct {
	Process config.ProcessJson
	Conn *net.Conn
	Request Request
}
// Structure représentant une requête d'un processus i
type Request struct {
	Message string
	IdProcess int
	Clock Clock
}


/* FUNCTIONS */

// Retourne le processRequest correspondant au processus local
func GetCurrentConnectionProcess() ConnectionProcess {
	return processRequest[idCurrentProcess]
}


// Incrémente l'estampille
func (c *Clock) Increment() uint64 {
	return atomic.AddUint64(&c.Timestamp,1)
}

// Mets la valeur de l'estampille à jour selon ce qui est reçu
func (c *Clock) UpdateClock(v Clock) {
	cur := atomic.LoadUint64(&c.Timestamp)

	if v.Timestamp >= cur {
		atomic.CompareAndSwapUint64(&c.Timestamp,cur,v.Timestamp + 1)
	}
}


// Fonction REQ de Lamport
func REQ(idProcessor int, clockProcessor Clock) {
	connectionProcessor := processRequest[idProcessor]
	clock.UpdateClock(clockProcessor)

	// Mise à jour de la commande / clock
	connectionProcessor.Request.Message = "REQ"
	connectionProcessor.Request.Clock = clockProcessor
	processRequest[connectionProcessor.Process.Id] = connectionProcessor

	// TO DO : broadcast message ACK

	CheckSC()
}

// Fonction ACK de Lamport
func ACK(idProcessor int, clockProcessor Clock) {
	connectionProcessor := processRequest[idProcessor]
	clock.UpdateClock(clockProcessor)
	if connectionProcessor.Request.Message != "REQ" {
		connectionProcessor.Request.Message = "ACK"
		connectionProcessor.Request.Clock = clockProcessor
		processRequest[idProcessor] = connectionProcessor
	}

	CheckSC()
}

// Fonction REL de Lamport
func REL(idProcessor int, clockProcessor Clock) {
	connectionProcessor := processRequest[idProcessor]
	clock.UpdateClock(clockProcessor)
	connectionProcessor.Request.Message = "REL"
	connectionProcessor.Request.Clock = clockProcessor
	processRequest[idProcessor] = connectionProcessor

	CheckSC()
}


// Teste si le processus courant à le droit d'accéder à la section critique
func CheckSC() {
	if GetCurrentConnectionProcess().Request.Message != "REQ" {
		return
	}

	older := true

	for _, connectionProcessor := range processRequest {
		if connectionProcessor.Process.Id != GetCurrentConnectionProcess().Process.Id {
			if connectionProcessor.Request.Clock.Timestamp < GetCurrentConnectionProcess().Request.Clock.Timestamp {
				older = false
				break
			} else if connectionProcessor.Request.Clock.Timestamp == GetCurrentConnectionProcess().Request.Clock.Timestamp &&
				connectionProcessor.Process.Id < GetCurrentConnectionProcess().Process.Id {
				older = false
				break
			}
		}
	}
	if older {
		common.Debug("The process can access to the SC now")
		canSc = true
	}
}

// Mets le processus en attente de la Section Critique
func WaitForSC() {
	for i := 0; !canSc; i++ {
		if i%100 == 0 {
			common.Debug("The process is waiting for the SC")
		}
	}

}

// Libère la Section Critique
func FinishSC() {
	clock.Increment()
	currentProc := GetCurrentConnectionProcess()
	currentProc.Request.Message = "REL"
	currentProc.Request.Clock = clock

	canSc = false
	common.Debug("The SC is finished")

	processRequest[currentProc.Process.Id] = currentProc
	for _, connProc :=range processRequest {
		if connProc.Process.Id != idCurrentProcess {
			// TO DO : broadcast message de REL
		}
	}
}
