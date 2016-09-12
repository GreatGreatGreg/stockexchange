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

class Portfolio extends React.Component {
  constructor() {
    super();
  }

  render() {
    if(!this.props.value.shares || this.props.value.shares.length == 0) {
      return (
        <div className="panel panel-info">
          <div className="panel-heading">
            <h3 className="panel-title">
              <span>Portfolio</span>
            </h3>
          </div>
          <div className="panel-body">
            <h4 className="text-center">Your portfolio is empty</h4>
          </div>
          <div className="panel-footer">
            <span className="text-right text-uppercase"><strong>Cash ${this.props.value.balance}</strong></span>
          </div>
        </div>
      );
    }

    return (
      <div className="panel panel-info">
        <div className="panel-heading">
          <h3 className="panel-title">
            <span>Portfolio</span>
          </h3>
        </div>
        <div className="panel-body">
          <table className="table table-hover">
            <thead>
              <tr>
                <th>Symbol</th>
                <th>Name</th>
                <th>Paid Price</th>
                <th>Quantity</th>
              </tr>
            </thead>
            <tbody>
              {
                this.props.value.shares.map(function(item) {
                  return (
                    <tr key={item.symbol}>
                      <th scope="row">{item.symbol}</th>
                      <td>{item.name}</td>
                      <td>{item.paidPrice}</td>
                      <td>{item.quantity}</td>
                    </tr>
                  );
                })
              }
            </tbody>
          </table>
        </div>
        <div className="panel-footer">
          <span className="text-right text-uppercase"><strong>Cash ${this.props.value.balance}</strong></span>
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
    if(this.props.message) {
      return (
        <div className="panel panel-danger">
          <div className="panel-body">
            <h4 className="text-center">{this.props.message}</h4>
          </div>
        </div>
      );
    }

    if(this.props.result && this.props.result.length > 0) {
      return (
        <div className="panel panel-warning">
          <div className="panel-heading">
            <h3 className="panel-title">
              <span>Stock</span>
            </h3>
          </div>
          <div className="panel-body">
            <table className="table table-hover">
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
                <input type="text" className="form-control" aria-label="hidden" ref="quantityBox" />
                <div className="input-group-btn">
                  <button type="button" className="btn btn-default" onClick={() => this.props.onBuyClick(this.props.result[0], this.refs.quantityBox.value)}>Buy</button>
                  <button type="button" className="btn btn-default" onClick={this.props.sell}>Sell</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return <div></div>;
  }
}

class Alert extends React.Component {
  constructor() {
    super()
    this.close = this.close.bind(this);
    this.state = {visible: true};
  }

  render() {
    if(this.state.visible && this.props.message) {
      return (
        <div className={"alert alert-master alert-" + this.props.level} role="alert">
          <button type="button" className="close" onClick={this.close}>×</button>
          <span>{this.props.message}</span>
        </div>
      );
    }
    return <div></div>;
  }

  close() {
    this.setState({visible: false});
  }
}

class ApplicationContainer extends React.Component {
  constructor() {
    super();

    this.state = {
      portfolio: {balance: 100000, shares:[]},
      search:{result: [], message: ""},
      alert: {level: "info", message: ""},
    };

    this.search = this.search.bind(this);
    this.buy = this.buy.bind(this);
    this.sell = this.sell.bind(this);
  }

  render() {
    return (
      <div>
        <NavigationBar onSearchClick={this.search} />
        <div className="container container-small">
          <Portfolio value={this.state.portfolio} />
          <SearchContainer result={this.state.search.result} message={this.state.search.message} onBuyClick={this.buy} onSellClick={this.sell}/>
          <Alert level={this.state.alert.level} message={this.state.alert.message}/>
        </div>
      </div>
    );
  }

  componentDidMount() {
    $.ajax({
      url: "/api/v1/portfolio",
      dataType: 'json',
      error: function() {
        this.setState({alert: {level: "danger", message:""}});
      }.bind(this),
      success: function(portfolio) {
        this.setState({portfolio: portfolio});
      }.bind(this)
    });
  }

  search(text) {
    $.ajax({
      url: "/api/v1/search?query="+text,
      dataType: 'json',
      error: function() {
        this.setState({portfolio: this.state.portfolio, search:{result: [], message: "Nothing has been found"}})
      }.bind(this),
      success: function(stock) {
        this.setState({portfolio: this.state.portfolio, search:{result: stock, message: ""}})
      }.bind(this)
    });
  }

  buy(stock, quantity) {
    $.ajax({
      url: "/api/v1/buy?quantity="+quantity,
      contentType: "application/json; charset=utf-8",
      type: "POST",
      dataType: 'json',
      data:  JSON.stringify(stock),
      error: function(jqXHR, textStatus, errorThrown) {
      this.setState({alert: {level: "danger", message: jqXHR.responseText}});
      }.bind(this),
      success: function(portfolio) {
        this.setState({portfolio: portfolio, search: this.state.search})
      }.bind(this)
    });
  }

  sell(symbol, price, quantity) {

  }
}

ReactDOM.render(<ApplicationContainer/>, document.getElementById('app'));
