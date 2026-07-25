import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Home from './pages/Home'
import Locations from './pages/Locations'
import LocationDetails from './pages/LocationDetails'
import Bookings from './pages/Bookings'
import Profile from './pages/Profile'
import Host from './pages/Host'

function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/locations" element={<Locations />} />
        <Route path="/locations/:id" element={<LocationDetails />} />
        <Route path="/bookings" element={<Bookings />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/host" element={<Host />} />
      </Routes>
    </Layout>
  )
}

export default App
