import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './components/Home';
import Calculator from './components/Calculator';
import './App.css';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <nav className="bg-white shadow-md p-4">
          <ul className="flex space-x-6 justify-center">
            <li>
              <Link to="/" className="text-blue-600 hover:text-blue-800 font-medium">Home</Link>
            </li>
            <li>
              <Link to="/calculator" className="text-blue-600 hover:text-blue-800 font-medium">Calculator</Link>
            </li>
          </ul>
        </nav>

        <div className="container mx-auto p-4">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/calculator" element={<Calculator />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
