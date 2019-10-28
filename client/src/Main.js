import React, { Component } from 'react';
import { Switch, Route } from 'react-router-dom';
import HomePage from './components/HomePage';

class Main extends Component {
   render() {
       return (
        <Switch>
            <Route exact path="/" component={HomePage} />
        </Switch>
       );
   }
}

export default Main;
