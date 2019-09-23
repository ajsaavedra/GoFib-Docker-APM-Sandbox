import React, { Component } from 'react';
import { Switch, Route } from 'react-router-dom';
import HomePage from './components/HomePage';
import AllValuesPage from './components/AllValuesPage';

class Main extends Component {
   render() {
       return (
        <Switch>
            <Route exact path="/" component={HomePage} />
            <Route exact path="/all-values" component={AllValuesPage} />
        </Switch>
       );
   }
}

export default Main;
