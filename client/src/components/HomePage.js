import React, { Component } from 'react';
import { Form, FormControl, Button } from 'react-bootstrap';
import TableValues from './TableValues';

class HomePage extends Component {
   render() {
       return (
           <div className="main-content">
               <h1>Fibonacci Calculator</h1>
               <div>
                <Form inline>
                    <FormControl type="text" placeholder="fib(n)" className="mr-sm-2" />
                    <Button variant="outline-dark">Calculate</Button>
                </Form>
                <TableValues/>
               </div>
            </div>
       );
   }
}

export default HomePage;
