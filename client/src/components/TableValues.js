import React, { Component } from 'react';
import { Table, Button } from 'react-bootstrap';

class TableValues extends Component {
    constructor() {
        super();
        this.state = {
            values: [
                {index: 1, value: 1, execution:'1s'},
                {index: 3, value: 3, execution:'1s'},
                {index: 5, value: 5, execution:'2s'}]
        };
    }

    renderTableValues() {
        return this.state.values.map(fib =>
            <tr key={fib.index}>
                <td>{fib.index}</td>
                <td>{fib.value}</td>
                <td>{fib.execution}</td>
                <td><Button variant="outline-danger">Delete</Button></td>
            </tr>
        )
    }

    render() {
       return (
           <Table striped bordered hover>
           <thead>
             <tr>
               <th>Index</th>
               <th>Value</th>
               <th>Execution Time</th>
               <th></th>
             </tr>
           </thead>
           <tbody>
             {this.renderTableValues()}
           </tbody>
         </Table>
       );
   }
}

export default TableValues;
