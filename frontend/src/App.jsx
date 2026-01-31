import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import Home from './pages/Home'
import Navigation from './components/Navigation/Navigation'
import Authenticate from './pages/Authenticate/Authenticate'
import GuestRoute from './routes/GuestRoute'
import SemiProtectedRoute from './routes/SemiProtectedRoute'
import Activate from './pages/Activate/Activate'

function App() {
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
        {/* <Route path="/" element={
          <GuestRoute>
            <Home />
          </GuestRoute>
        } />
        <Route path="/authenticate" element={
          <GuestRoute>
            <Authenticate />
          </GuestRoute>
        } /> */}
      </Routes>
    </Router>
  )
} 

export default App
