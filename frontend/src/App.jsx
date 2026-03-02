import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import Home from './pages/Home'
import Navigation from './components/Navigation/Navigation'
import Authenticate from './pages/Authenticate/Authenticate'
import GuestRoute from './routes/GuestRoute'
import SemiProtectedRoute from './routes/SemiProtectedRoute'
import Activate from './pages/Activate/Activate'
import ProtectedRoute from './routes/ProtectedRoute'
import Rooms from './pages/Rooms/Rooms'
import { useLoadingWithRefresh } from './hooks/useLoadingWithRefresh'
import Loader from './components/Loader/Loader'
import Room from './pages/Room/Room'

function App() {
  const { loading } = useLoadingWithRefresh();


  if (loading) {
    return (
      <Loader message="Please wait..." />
    )
  }
  return (
    <Router>
      <Navigation />
      <Routes>
        <Route element={<GuestRoute />}>
          <Route path="/" element={<Home />} />
          <Route path="/authenticate" element={<Authenticate />} />
        </Route>
        <Route element={<SemiProtectedRoute />}>
          <Route path="/activate" element={<Activate />} />
        </Route>
        <Route element={<ProtectedRoute />}>
          <Route path="/rooms" element={<Rooms />} />
        </Route>
        <Route element={<ProtectedRoute />}>
          <Route path="/room/:id" element={<Room />} />
        </Route>
      </Routes>
    </Router>
  )
}

export default App
