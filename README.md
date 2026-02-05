# 🎓 Projet Middleware - Système d'Alertes EDT

Système de notification automatique des changements d'emploi du temps de l'Université Clermont Auvergne (UCA).

## 📋 Description

Ce projet implémente une architecture microservices permettant de :
- Récupérer les événements de l'emploi du temps UCA
- Détecter les modifications (changement de salle, horaire, etc.)
- Envoyer des alertes par email aux utilisateurs abonnés

## 🏗️ Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Scheduler  │────▶│    NATS     │────▶│  Timetable  │────▶│   Alerter   │
│             │     │  JetStream  │     │   Consumer  │     │   Consumer  │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
      │                                        │                    │
      ▼                                        ▼                    ▼
┌─────────────┐                         ┌─────────────┐     ┌─────────────┐
│  UCA EDT    │                         │   SQLite    │     │  Mail API   │
│    API      │                         │  (events)   │     │             │
└─────────────┘                         └─────────────┘     └─────────────┘
      │
      ▼
┌─────────────┐
│ Config API  │
│  (agendas,  │
│   alerts)   │
└─────────────┘
```

## 🧩 Composants

### 1. Config API (`/config`)
API REST pour gérer les agendas et les alertes.

**Endpoints :**
- `GET /agendas` - Liste des agendas
- `POST /agendas` - Créer un agenda
- `GET /agendas/{id}` - Détail d'un agenda
- `PUT /agendas/{id}` - Modifier un agenda
- `DELETE /agendas/{id}` - Supprimer un agenda
- `GET /alerts` - Liste des alertes
- `POST /alerts` - Créer une alerte
- `GET /alerts/{id}` - Détail d'une alerte
- `PUT /alerts/{id}` - Modifier une alerte
- `DELETE /alerts/{id}` - Supprimer une alerte

### 2. Scheduler (`/scheduler`)
Récupère périodiquement les événements depuis l'API UCA et les publie sur NATS.

**Fonctionnalités :**
- Récupération des agendas depuis Config API
- Parsing des fichiers iCal de l'UCA
- Publication des événements sur le stream `EVENTS`

### 3. Timetable (`/timetable`)
Consumer NATS qui stocke les événements et détecte les changements.

**Fonctionnalités :**
- Consommation du stream `EVENTS`
- Stockage des événements en SQLite
- Détection des modifications (salle, horaire, etc.)
- Publication des alertes sur le stream `ALERTS`

### 4. Alerter (`/alerter`)
Consumer NATS qui envoie les notifications par email.

**Fonctionnalités :**
- Consommation du stream `ALERTS`
- Récupération des alertes configurées depuis Config API
- Envoi d'emails via l'API mail-api.edu.forestier.re

## 🚀 Installation

### Prérequis
- Go 1.21+
- Docker (pour NATS)
- Compte UCA (pour le token mail)

### 1. Cloner le projet
```bash
git clone https://github.com/SaherNadine/Projet-Middleware.git
cd Projet-Middleware
```

### 2. Lancer NATS avec JetStream
```bash
docker run -d --name nats -p 4222:4222 nats -js
```

### 3. Installer les dépendances
```bash
# Config API
cd config && go mod tidy && cd ..

# Scheduler
cd scheduler && go mod tidy && cd ..

# Timetable
cd timetable && go mod tidy && cd ..

# Alerter
cd alerter && go mod tidy && cd ..
```

## 🎮 Utilisation

### 1. Démarrer Config API
```bash
cd config
go run cmd/main.go
# ✅ API Config started on :8080
```

### 2. Démarrer Timetable
```bash
cd timetable
go run cmd/main.go --debug
# ✅ Consumer Timetable en écoute sur EVENTS.>
```

<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/52410342-a068-4f81-abe7-bb85ee09eeb5" />


### 3. Démarrer Alerter
```bash
cd alerter
go run cmd/main.go \
  --mail-api="https://mail-api.edu.forestier.re" \
  --mail-token="VOTRE_TOKEN" \
  --debug
# ✅ Consumer Alerter en écoute...
```

<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/7debec54-8873-4e06-be13-d8a68054bdc6" />



<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/d985ca2c-586a-4de2-8bd3-3c244ba18f10" />


### 4. Démarrer Scheduler
```bash
cd scheduler
go run cmd/main.go
# ✅ Scheduler en cours d'exécution
```
<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/c4ed2b7d-8339-4d26-bb9c-2d5b5c99941d" />

## ⚙️ Configuration

### Créer des agendas
```bash
# M1 Informatique - Groupe 1
curl -X POST http://localhost:8080/agendas \
  -H "Content-Type: application/json" \
  -d '{"name": "M1 Groupe 1 langue", "uca_id": "13295"}'

# M1 Informatique - Groupe 2
curl -X POST http://localhost:8080/agendas \
  -H "Content-Type: application/json" \
  -d '{"name": "M1 Groupe 2 langue", "uca_id": "13345"}'
```

### Créer des alertes
```bash
# Alerte pour tous les événements
curl -X POST http://localhost:8080/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "recipient": "votre.email@etu.uca.fr",
    "condition": "always",
    "is_active": true
  }'

# Alerte uniquement pour les changements de salle
curl -X POST http://localhost:8080/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "recipient": "votre.email@etu.uca.fr",
    "condition": "room_change",
    "is_active": true
  }'
```



### Conditions d'alerte disponibles
| Condition | Description |
|-----------|-------------|
| `always` | Toutes les modifications |
| `room_change` | Changement de salle uniquement |
| `time_change` | Changement d'horaire uniquement |
| `new_event` | Nouveaux événements uniquement |
| `deleted` | Événements supprimés uniquement |

## 🔑 Obtenir un token mail

1. Aller sur https://mail.edu.forestier.re
2. Entrer votre email `@etu.uca.fr`
3. Valider avec le code reçu par email
4. Copier le token généré

## 📁 Structure du projet

```
middleware/
├── config/                 # API Config
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── controllers/
│   │   │   ├── agendas/
│   │   │   └── alerts/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── services/
│   │   └── helpers/
│   └── go.mod
│
├── scheduler/              # Scheduler
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── models/
│   │   └── services/
│   └── go.mod
│
├── timetable/              # Timetable Consumer
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── controllers/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── services/
│   │   └── helpers/
│   └── go.mod
│
├── alerter/                # Alerter Consumer
│   ├── cmd/main.go
│   ├── config/
│   │   ├── template.go
│   │   └── templates/
│   ├── internal/
│   │   ├── helpers/
│   │   ├── models/
│   │   └── services/
│   └── go.mod
│
└── README.md
```

## 🔄 Flux de données

1. **Scheduler** récupère les agendas depuis Config API
2. **Scheduler** fetch les événements iCal depuis l'UCA
3. **Scheduler** publie les événements sur `EVENTS.created`
4. **Timetable** consomme les événements
5. **Timetable** compare avec la BDD et détecte les changements
6. **Timetable** publie les alertes sur `ALERTS.created` ou `ALERTS.modified`
7. **Alerter** consomme les alertes
8. **Alerter** récupère les destinataires depuis Config API
9. **Alerter** envoie les emails via l'API mail

## 🛠️ Technologies

- **Go** - Langage de programmation
- **NATS JetStream** - Message broker
- **SQLite** - Base de données
- **Chi** - Router HTTP
- **Logrus** - Logging

## 👥 Équipe

- Ibrahim EL HAOURARI
- Saher Nadine
- Aya

## 📝 Licence

Projet réalisé dans le cadre du cours de Middleware - ISIMA - Université Clermont Auvergne

---

## 🐛 Dépannage

### NATS ne démarre pas
```bash
docker rm -f nats
docker run -d --name nats -p 4222:4222 nats -js
```

### Erreur "invalid token" pour l'API mail
- Vérifiez que vous utilisez `https://mail-api.edu.forestier.re` (avec `-api`)
- Régénérez un token sur https://mail.edu.forestier.re

### Les emails ne sont pas reçus
- Vérifiez le dossier spam
- Vérifiez que l'adresse est bien `@etu.uca.fr`

### "no such table: events"
```bash
rm -f timetable/events.db
# Relancer timetable
```

### "no such table: alerts"
```bash
rm -f config/config.db
# Relancer config
```
