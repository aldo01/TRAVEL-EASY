# MongoDB Setup for Travel Easy

You need MongoDB to run the backend. Choose one option:

## Option 1: MongoDB Atlas (Cloud - Recommended for Windows)

**Easiest option - No local installation needed!**

1. **Create Free Account:**
   - Go to https://www.mongodb.com/cloud/atlas/register
   - Sign up with email or Google

2. **Create Cluster:**
   - Click "Build a Database"
   - Choose "FREE" tier (M0 Sandbox)
   - Select a cloud provider and region (choose closest to you)
   - Click "Create Cluster" (takes 3-5 minutes)

3. **Configure Database Access:**
   - Go to "Database Access" in left menu
   - Click "Add New Database User"
   - Create username and password (save these!)
   - User Privileges: "Read and write to any database"
   - Click "Add User"

4. **Configure Network Access:**
   - Go to "Network Access" in left menu
   - Click "Add IP Address"
   - Click "Allow Access from Anywhere" (for development)
   - Click "Confirm"

5. **Get Connection String:**
   - Go to "Database" → Click "Connect"
   - Choose "Connect your application"
   - Copy the connection string (looks like):
     ```
     mongodb+srv://<username>:<password>@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
     ```

6. **Update Backend .env:**
   ```env
   MONGODB_URI=mongodb+srv://youruser:yourpass@cluster0.xxxxx.mongodb.net/travel-easy?retryWrites=true&w=majority
   ```
   Replace `<username>` and `<password>` with your credentials

## Option 2: Local MongoDB (Windows)

1. **Download:**
   - Go to https://www.mongodb.com/try/download/community
   - Select "Windows" and download MSI installer

2. **Install:**
   - Run the MSI installer
   - Choose "Complete" installation
   - Install "MongoDB as a Service" (check this option)
   - Install MongoDB Compass (GUI tool)

3. **Start MongoDB:**
   - MongoDB should auto-start as a Windows service
   - Or manually:
     ```bash
     net start MongoDB
     ```

4. **Verify Installation:**
   ```bash
   mongod --version
   ```

5. **Backend .env:**
   ```env
   MONGODB_URI=mongodb://localhost:27017/travel-easy
   ```

## After Setup

### Start Backend Server
```bash
cd backend
npm run dev
```

### Test Connection
Open http://localhost:5000/health

You should see:
```json
{
  "success": true,
  "message": "Travel Easy API is running"
}
```

### Seed Database (Optional)
```bash
npx ts-node src/utils/seed.ts
```

## Troubleshooting

### MongoDB Atlas Connection Issues:
- Check username/password are correct
- Ensure IP is whitelisted (0.0.0.0/0 for development)
- Wait a few minutes after creating cluster

### Local MongoDB Issues:
- Check if service is running: `net start MongoDB`
- Check port 27017 is not blocked
- Try restarting: `net stop MongoDB` then `net start MongoDB`

## Recommended: Use MongoDB Atlas

For Windows development, **MongoDB Atlas is easier** because:
- ✅ No local installation
- ✅ Works from anywhere
- ✅ Free tier available
- ✅ Automatic backups
- ✅ No configuration needed

---

Once MongoDB is set up, you're ready to run the full application!
