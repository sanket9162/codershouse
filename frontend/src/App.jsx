import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import Home from './pages/Home'
import Navigation from './components/Navigation/Navigation'
import Authenticate from './pages/authenticate/Authenticate'
import GuestRoute from './routes/GuestRoute'


function App() {
  return (
    <Router>
      <Navigation />
      <Routes>

        <Route element={<GuestRoute />}>
    <Route path="/" element={<Home />} />
    <Route path="/authenticate" element={<Authenticate />} />
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
