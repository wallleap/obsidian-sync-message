# OB Sync

A three-part synchronization system for Obsidian notes:
- **Server**: Go + SQLite backend
- **Frontend**: React web interface for sending messages/attachments
- **Obsidian Plugin**: Sync notes from server to your vault

## Server

### Requirements
- Go 1.22+

### Installation
```bash
cd server
go mod download
```

### Running
```bash
cd server
go run cmd/main.go
```

The server will start on `http://localhost:8080`

### API Endpoints
- `POST /api/user/generate` - Generate a new user ID
- `POST /api/user/validate` - Validate a user ID
- `POST /api/message/send` - Send a text/URL message
- `POST /api/message/upload` - Upload an attachment
- `POST /api/message/sync` - Sync messages for a user
- `GET /api/message/file/:id` - Download a file attachment

## Frontend

### Requirements
- Node.js 18+
- npm

### Installation
```bash
cd frontend
npm install
```

### Development
```bash
cd frontend
npm run dev
```

### Building
```bash
cd frontend
npm run build
```

## Obsidian Plugin

### Installation
1. Build the plugin:
   ```bash
   cd obsidian-sample-plugin
   npm install
   npm run build
   ```

2. Copy `main.js`, `manifest.json`, and `styles.css` to your Obsidian vault's plugin folder:
   ```
   <Vault>/.obsidian/plugins/ob-sync/
   ```

3. Enable the plugin in Obsidian settings.

### Configuration
- **User ID**: Your unique identifier for synchronization
- **Server URL**: The URL of your OB Sync server
- **Save Folder**: Folder to store synced notes
- **Attachment Folder**: Subfolder for attachments
- **Image Folder**: Subfolder for images
- **Time Format**: Format for timestamps
- **Title Template**: Template for note titles
- **Use Image Bed Relay**: Enable image upload to image bed

## Usage

1. **Generate User ID**: On the frontend, generate a new user ID and save it
2. **Configure Plugin**: Set your user ID and server URL in the plugin settings
3. **Send Messages**: Use the frontend to send text, URLs, or attachments
4. **Sync**: Click the plugin icon or use the command to sync messages to your vault

## Project Structure
```
ob-sync/
├── server/                 # Go backend
│   ├── cmd/               # Entry point
│   ├── config/            # Configuration
│   ├── internal/          # Internal packages
│   │   ├── handler/       # HTTP handlers
│   │   ├── model/         # Database models
│   │   ├── repository/    # Data access
│   │   └── util/          # Utilities (logger)
│   └── go.mod
├── frontend/              # React frontend
│   ├── src/
│   │   ├── api/          # API client
│   │   ├── components/   # React components
│   │   └── hooks/        # Custom hooks
│   └── package.json
└── obsidian-plugin/ # Obsidian plugin
    ├── src/
    │   ├── main.ts       # Plugin entry
    │   └── settings.ts   # Settings
    └── package.json
```
