import React from 'react';
import {
  ChakraProvider,
  theme,
} from '@chakra-ui/react';
import {
  BrowserRouter as Router,
  Route,
  Switch,
} from 'react-router-dom';
import { HomePage, TermsPage } from './pages';

export default function App() {
  return (
    <ChakraProvider theme={theme}>
      <Router>
        <Switch>
          <Route path='/terms'>
            <TermsPage />
          </Route>
          <Route path='/'>
            <HomePage />
          </Route>
        </Switch>
      </Router>
    </ChakraProvider>
  );
}
