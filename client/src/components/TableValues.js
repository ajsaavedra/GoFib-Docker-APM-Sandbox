import React, { Component } from 'react';
import { Table, Button } from 'react-bootstrap';
import { connect } from 'react-redux';
import * as actions from '../actions';

class TableValues extends Component {
    constructor() {
      super();
      this.state = {
        values: null
      }
      this.renderTableValues = this.renderTableValues.bind(this);
    }

    componentWillMount() {
        this.props.getFibValues();
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.values) {
          this.setState({
            values: nextProps.values
          });
        }
        if (nextProps.deleted) {
          let update = [...this.state.values];
          update = update.filter(values => values.idx != nextProps.value);
          this.setState({ values: update });
        }
    }

    handleDelete(index) {
      this.props.deleteFib(index);
    }

    renderTableValues() {
        if (this.state.values) {
            return this.state.values.map(fib =>
                <tr key={fib.idx}>
                    <td>{fib.idx}</td>
                    <td>{fib.fib}</td>
                    <td>{fib.elapsed}</td>
                    <td><Button variant="outline-danger" onClick={() => this.handleDelete(fib.idx)}>Delete</Button></td>
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
    const { error, values, deleted, value } = state.fib;
    return { error, values, deleted, value };
}

export default connect(mapStateToProps, actions)(TableValues);
