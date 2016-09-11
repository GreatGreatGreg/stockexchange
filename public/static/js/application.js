class NavigationBar extends React.Component {
  constructor() {
    super();
  }

  render() {
    return (
      <nav className="navbar navbar-inverse navbar-static-top">
        <div className="container">
          <div className="navbar-header">
            <a className="navbar-brand" href="#">StockExchange</a>
          </div>
          <div id="navbar" className="collapse navbar-collapse">
            <div className="navbar-form navbar-right">
              <div className="input-group">
                <input type="text" ref="searchBox" className="form-control" placeholder="Search for..."/>
                <span className="input-group-btn">
                  <button className="btn btn-default" type="button" onClick={() => this.props.onSearchClick(this.refs.searchBox.value)}>
                    <span className="glyphicon glyphicon-search" aria-hidden="true"></span>
                  </button>
                </span>
              </div>
            </div>
          </div>
        </div>
      </nav>
    );
  }
}

class GridView extends React.Component {
  constructor() {
    super();
  }

  render() {
    let element;

    if(this.props.source.length == 0) {
      element = (<h4 className="text-center">{this.props.noContentText}</h4>);
    } else {
      element = (
        <table className="table">
          <thead>
            <tr>
              <th>Symbol</th>
              <th>Name</th>
              <th>Ask Price</th>
              <th>Bid Price</th>
            </tr>
          </thead>
          <tbody>
            {
              this.props.source.map(function(item) {
                return (
                  <tr key={item.symbol}>
                    <th scope="row">{item.symbol}</th>
                    <td>{item.name}</td>
                    <td>{item.askPrice}</td>
                    <td>{item.bidPrice}</td>
                  </tr>
                );
              })
            }
          </tbody>
        </table>
      );
    }

    return element
  }
}

class Portfolio extends React.Component {
  constructor() {
    super();
    this.state = {balance: 100000, stock: []}
  }

  render() {
    return (
      <div className="panel panel-info">
        <div className="panel-heading">
          <h3 className="panel-title">
            <span>Portfolio</span>
            <span className="text-right pull-right">${this.state.balance}</span>
          </h3>
        </div>
        <div className="panel-body">
          <GridView source={this.state.stock} noContentText="Your portfolio is empty"/>
        </div>
      </div>
    );
  }
}

class SearchContainer extends React.Component {
  constructor() {
    super();
  }

  render() {
    if(this.props.infoMessage) {
      return (
        <div className="panel panel-info">
          <div className="panel-body">
            <h4 className="text-center">{this.props.infoMessage}</h4>
          </div>
        </div>
      );
    }

    if(this.props.result.length > 0) {
      let grid = (
        <div className="panel panel-info">
          <div className="panel-body">
            <table className="table">
              <thead>
                <tr>
                  <th>Symbol</th>
                  <th>Name</th>
                  <th>Ask Price</th>
                  <th>Bid Price</th>
                </tr>
              </thead>
              <tbody>
                {
                  this.props.result.map(function(item) {
                    return (
                      <tr key={item.symbol}>
                        <th scope="row">{item.symbol}</th>
                        <td>{item.name}</td>
                        <td>${item.askPrice}</td>
                        <td>${item.bidPrice}</td>
                      </tr>
                    );
                  })
                }
              </tbody>
            </table>
          </div>
          <div className="panel-footer">
            <div className="input-group">
              <div className="input-group">
                <span className="input-group-addon">Quantity</span>
                <input type="text" className="form-control" aria-label="hidden"/>
                <div className="input-group-btn">
                  <button type="button" className="btn btn-default">Buy</button>
                  <button type="button" className="btn btn-default">Sell</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      );

      if(this.props.errorMessage) {
        return (
          <div>
            {grid}
            <div className="alert alert-danger" role="alert">{this.props.errorMessage}</div>
          </div>
        );
      } else {
        return grid;
      }
    }

    return <div/>;
  }
}

class ApplicationContainer extends React.Component {
  constructor() {
    super();
    this.state = {result: [], message: ""}
    this.search = this.search.bind(this);
  }

  render() {
    return (
      <div>
        <NavigationBar onSearchClick={this.search} />
        <div className="container container-small">
          <SearchContainer result={this.state.result} infoMessage={this.state.message} />
        </div>
      </div>
    );
  }

  search(text) {
    $.ajax({
      url: "/api/v1/search?query="+text,
      dataType: 'json',
      error: function() {
        this.setState({message: "Nothing has been found", result: []});
      }.bind(this),
      success: function(shares) {
        this.setState({message: "", result: shares});
      }.bind(this)
    });
  }
}

ReactDOM.render(<ApplicationContainer/>, document.getElementById('app'));

