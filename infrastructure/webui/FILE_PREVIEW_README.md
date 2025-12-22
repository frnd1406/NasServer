# File Preview System - Benutzerhandbuch

## 🎯 Überblick

Das **File Preview System** ist eine intelligente Dateiauswahl- und Vorschau-Lösung für NAS.AI. Es kombiniert:

- **Semantic Search** mit AI-gestützter Relevanz-Bewertung
- **Intelligente Dateiauswahl** (Auto-open bei 1 Datei, Selection bei mehreren)
- **Multi-Format Viewer** (txt, json, code, images, markdown, etc.)
- **Interaktive Chat-Integration** mit RAG (Retrieval Augmented Generation)

---

## 🚀 Features

### 1. **Intelligente Dateiauswahl**
- **Single File**: Automatisches Öffnen wenn nur eine Datei relevant ist
- **Multiple Files**: Interaktive Auswahl mit Similarity-Scores
- **Relevanz-Anzeige**: Visuelle Darstellung der AI-basierten Relevanz (0-100%)

### 2. **Multi-Format File Viewer**
Unterstützte Dateitypen:
- **Text**: `.txt`, `.log`, `.md`
- **Code**: `.js`, `.jsx`, `.ts`, `.tsx`, `.py`, `.go`
- **Strukturiert**: `.json`, `.xml`, `.csv`
- **Bilder**: `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`
- **Dokumente**: `.pdf` (geplant)

### 3. **Enhanced Chat Widget**
- Natürliche Sprach-Anfragen
- Automatische Quellenangaben
- Direkte File-Preview aus Chat
- Navigation zwischen mehreren Dateien

---

## 📁 Komponenten-Architektur

```
webui/src/
├── components/
│   ├── FilePreviewPanel.jsx      # Datei-Vorschau rechts
│   ├── FileSelector.jsx          # Interaktive Dateiauswahl
│   └── EnhancedChatWidget.jsx    # Chat mit File-Integration
├── pages/
│   └── FilePreviewDemo.jsx       # Demo-Seite zum Testen
└── Layout.jsx                    # Hauptlayout mit Chat
```

### Backend (Go API)
```
api/src/handlers/
└── file_content.go               # File Content Endpoint
```

---

## 🎮 Nutzung

### Demo-Seite öffnen
Besuche: **http://localhost/demo**

Die Demo-Seite zeigt 4 Szenarien:

1. **Single File Demo** - Eine Datei (Steuerbescheid)
2. **Multiple Files Demo** - Mehrere Finanz-Dokumente
3. **JSON File Demo** - Strukturierte API-Konfiguration
4. **Code File Demo** - Go Source Code

### Chat-Integration nutzen

1. **Chat öffnen**: Klicke auf das Chat-Icon rechts unten
2. **Frage stellen**: z.B. "Was kostet der Server?"
3. **Datei auswählen**: Wenn mehrere Dateien gefunden werden
4. **Vorschau öffnet sich**: Rechte Seite zeigt Datei-Inhalt

### Beispiel-Queries

```
"Zeige mir meine Steuerunterlagen"
"Wie viel habe ich mit Krypto verdient?"
"Was sind meine monatlichen Ausgaben?"
"Welche Rechnungen habe ich im März?"
```

---

## 🔧 API Endpoints

### File Content Abruf
```http
GET /api/v1/files/content?path=/mnt/data/example.txt
```

**Response:**
- Content-Type: Automatisch erkannt (text/plain, application/json, etc.)
- Body: Roher Datei-Inhalt
- Max Size: 10MB

**Security:**
- Nur Zugriff auf `/mnt/data/*`
- Keine Directory Traversal
- Validierung aller Pfade

### RAG Endpoint (existing)
```http
GET /api/v1/ask?q=Was kostet der Server?
```

**Response:**
```json
{
  "answer": "Der Server kostet 149,99€ [rechnung_xyz.txt]",
  "cited_sources": [
    {
      "file_id": "rechnung_xyz.txt",
      "file_path": "/mnt/data/finanzen/rechnung_xyz.txt",
      "similarity": 0.92
    }
  ],
  "confidence": "HOCH"
}
```

---

## 🎨 Komponenten-Details

### FilePreviewPanel

**Props:**
```javascript
{
  files: [              // Array von Dateien
    {
      file_path: string,
      file_id: string,
      content: string,
      similarity: number  // 0.0 - 1.0
    }
  ],
  currentIndex: number,  // Aktueller Index
  onClose: function,     // Callback zum Schließen
  onNavigate: function   // Callback für Navigation
}
```

**Features:**
- Automatische Typ-Erkennung
- Syntax-Highlighting für Code
- Zoom für Bilder (25% - 200%)
- Download-Funktion
- Navigation bei mehreren Dateien

### FileSelector

**Props:**
```javascript
{
  files: Array,           // Dateien zur Auswahl
  onSelectFile: function, // Callback bei Auswahl
  autoSelectSingle: bool  // Auto-open bei 1 Datei (default: true)
}
```

**Features:**
- Similarity-Score Visualisierung
- Farbcodierung (grün=hoch, blau=mittel, gelb=niedrig)
- Fortschrittsbalken für Relevanz
- Hover-Effekte

### EnhancedChatWidget

**Features:**
- RAG-Integration
- Automatische Datei-Erkennung
- File Selector bei mehreren Quellen
- Auto-open bei einzelner Quelle
- Preview-Panel Management

---

## 🧪 Testing

### Lokales Testing

1. **Start Docker Environment:**
```bash
cd infrastructure
docker-compose up -d
```

2. **Access Demo Page:**
```
http://localhost/demo
```

3. **Test Scenarios:**
- Klicke auf "Single File Demo" → Datei öffnet automatisch
- Klicke auf "Multiple Files Demo" → Wähle Datei aus
- Teste Chat mit echten Daten

### Test-Daten generieren

```bash
cd infrastructure/ai_knowledge_agent
python generate_corpus.py --all --output /mnt/data
```

Dies generiert realistische deutsche Dokumente:
- Steuerunterlagen
- Krypto-Auszüge
- Rechnungen
- Kontoauszüge

---

## 🔒 Security

### Pfad-Validierung
- Alle Pfade müssen mit `/mnt/data/` beginnen
- `filepath.Clean()` verhindert Directory Traversal
- Keine symbolischen Links erlaubt

### Size Limits
- Max Preview Size: **10MB**
- Größere Dateien werden abgelehnt

### Content-Type Detection
- Automatische Typ-Erkennung
- Override für bekannte Formate
- `X-Content-Type-Options: nosniff` Header

---

## 📊 Performance

### Optimierungen
- Content-Caching in Frontend
- Lazy Loading für große Dateien
- Parallele File-Content-Requests
- Debounced Search Queries

### Limits
- Max 10 Dateien in File Selector
- Max 10MB Datei-Größe
- Timeout: 8s für AI Requests

---

## 🚧 Roadmap

### Geplante Features

- [ ] PDF Viewer Integration
- [ ] Syntax Highlighting für mehr Sprachen
- [ ] File Thumbnails
- [ ] Keyboard Shortcuts (←/→ für Navigation)
- [ ] Dark/Light Mode für Code
- [ ] Copy-to-Clipboard für Code
- [ ] File Annotations
- [ ] Multi-File Compare View

---

## 🐛 Troubleshooting

### Datei öffnet nicht
- **Prüfen**: Ist Pfad korrekt? (`/mnt/data/*`)
- **Logs checken**: Browser Console + API Logs
- **File Size**: Größer als 10MB?

### Chat zeigt keine Quellen
- **AI Service**: Läuft Ollama?
- **Database**: Sind Embeddings indexiert?
- **Query**: Zu unspezifisch?

### Similarity Score niedrig
- **Embeddings**: Re-indexing notwendig?
- **Query**: Präziser formulieren
- **Model**: Ist mxbai-embed-large geladen?

---

## 📚 Weitere Ressourcen

- **API Docs**: http://localhost/swagger
- **Main Repo**: `/home/frnd14/f1406`
- **AI Agent**: `infrastructure/ai_knowledge_agent/`

---

## 👨‍💻 Development

### Neue Dateitypen hinzufügen

**FilePreviewPanel.jsx:**
```javascript
case 'pdf':
  return <PDFViewer content={fileContent} />;
```

**file_content.go:**
```go
case ".pdf":
  contentType = "application/pdf"
```

### Styling anpassen

Alle Komponenten nutzen **Tailwind CSS**:
```javascript
className="bg-slate-800 border border-white/10 rounded-xl"
```

---

## ✅ Checkliste für Deployment

- [ ] Environment Variables gesetzt
- [ ] Ollama Models geladen (`mxbai-embed-large`, `llama3.2`)
- [ ] Database Migrations durchgeführt
- [ ] Test-Daten generiert
- [ ] API Health Check OK (`/health`)
- [ ] Frontend Build erstellt (`npm run build`)
- [ ] CORS korrekt konfiguriert

---

**Version:** 1.0.0
**Datum:** 2025-12-07
**Entwickler:** NAS.AI Team
