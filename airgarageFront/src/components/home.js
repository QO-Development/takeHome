import React from 'react';
import { Navbar, Nav, Form, FormControl, Button, ListGroup } from 'react-bootstrap';

class Home extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      searchText: "",
      businesses: [],
    }
  }

  textChanged = (event) => {
    this.setState({
      searchText: event.target.value
    });
  }

  compare = (a, b) => {
    if (a.rating > b.rating) return 1;
    if (b.rating > a.rating) return -1;
  
    return 0;
  }

  buttonClicked = () => {

    fetch(`http://localhost:8080/api?searchText=${encodeURIComponent(this.state.searchText)}`, {
      method: 'GET',
    })
      .then((response) => {
        return response.json();
      })
      .then((json) => {
        console.log(json);

      
        
        json.businesses.sort(this.compare);

        this.setState({
          businesses: json.businesses,
        })
      })
      .catch((error) => {
        console.log("There was an error: ");
        console.log(error);
      });
  }

  render() {
    return (
      <>
        <Navbar bg="primary" variant="dark">
          <Navbar.Brand href="#home">AirGarage Homework</Navbar.Brand>
          <Nav className="mr-auto">
            <Nav.Link href="#home">Home</Nav.Link>
            <Nav.Link href="#features">Features</Nav.Link>
            <Nav.Link href="#pricing">Pricing</Nav.Link>
          </Nav>
          <Form inline>
            <FormControl onChange={this.textChanged} type="text" placeholder="Search" className="mr-sm-2" />
            <Button onClick={this.buttonClicked} variant="outline-light">Search</Button>
          </Form>
        </Navbar>

        <ListGroup>

          { this.state.businesses.map( ( value, index ) => {

            const score = ( value.review_count * value.rating ) / (value.review_count + 1);


            let img; 
            if (value.image_url) {
              img = <img alt="parking" src={value.image_url} style={{width: 100, height: 'auto'}} />;
            }

            return (
              <ListGroup.Item>
                <p>Address: {value.location.address1}</p>
                {img}
                <p>Star Rating: {value.rating} </p>
                <p>Review Count: {value.review_count} </p>
                <p> <a href={value.url}> Yelp Link </a> </p>
                <p>Parking Lot Score: {score} </p>

              </ListGroup.Item>

            );

          }) }

          
          
        </ListGroup>

      </>
    );
  }
}


export default Home;