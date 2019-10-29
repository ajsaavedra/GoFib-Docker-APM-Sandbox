import React, { Component } from 'react';
import { Table, Button } from 'react-bootstrap';
import { connect } from 'react-redux';
import * as actions from '../actions';

class TableValues extends Component {
    componentWillMount() {
        this.props.getFibValues();
    }

    renderTableValues() {
        if (this.props.values) {
            return this.props.values.map(fib =>
                <tr key={fib.idx}>
                    <td>{fib.idx}</td>
                    <td>{fib.fib}</td>
                    <td></td>
                    <td><Button variant="outline-danger">Delete</Button></td>
                </tr>
            )
        }
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

function mapStateToProps(state) {
    const { error, values } = state.fib;
    return { error, values };
}

export default connect(mapStateToProps, actions)(TableValues);
