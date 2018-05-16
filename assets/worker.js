var pages = [];
const orchestratorPage = require('webpage').create();

const EventEmitter = require('events');
const eventEmitter = new EventEmitter();

eventEmitter.on('next', (worker, page, link, kind) => {
  console.log('visiting: ' + link);

  page.open(link).then(function(status) {
    if (status === 'success') {
      sendMessageToOrchestrator('errormsg', {error: page.title}),
      if (page.title === 'Access Denied') {
        sendMessageToOrchestrator('errormsg', {error: 'Access Denied'});
      	slimer.exit();
      } else {
        switch (kind) {
          case 'categories':
            sendMessageToOrchestrator('categories', {
              worker: worker,
              categories: extractCategories(page),
            });
            break;
          case 'productsLinks':
            sendMessageToOrchestrator('productsLinks', {
              worker: worker,
              productsLinks: extractProducts(page),
              nextPage: extractNextPage(page),
            });
            break;
          case 'product':
            sendMessageToOrchestrator('product', {
              worker: worker,
              product: {...extractProduct(page), url: link},
            });
            break;
        }
      }
    }
  });
});

function processMessageFromOrchestrator(message) {
  if (message.workers) {
    for (let i = 0; i < message.workers; i++) {
      pages[i] = require('webpage').create();
    }
  } else {
    switch (message.command) {
      case 'extract':
        eventEmitter.emit('next', message.worker, pages[message.worker],
            message.link, message.kind);
    }
  }
}

function sendMessageToOrchestrator(type, content) {
  orchestratorPage.evaluate(
      (type, content) => {
        sendMessageToOrchestrator(type, content);
      },
      type,
      content,
  );
}

function connectToOrchestrator() {
  orchestratorPage.onConsoleMessage = function(message) {
    try {
      console.log(message);
      message = JSON.parse(message);

      processMessageFromOrchestrator(message);
    } catch (error) {
      console.error(error);
    }
  };

  function openPage() {
    orchestratorPage.open('http://localhost:4242/worker.html').
        then(function(status) {
          if (status === 'success') {
            console.log('loaded');
          } else {
            setTimeout(openPage, 25);
          }
        });
  }

  openPage();
}

connectToOrchestrator();

function extractCategories(page) {
  let results = page.evaluate(() => {
    return Array.from(
        document.querySelectorAll(
            '.navbar .navbar-inner .nav-collapse .nav li:first-child .cat li.dropdown .dropdown-menu .grid-demo h3:first-child a',
        ),
    );
  });

  return results.map(r => r.href);
}

function extractNextPage(page) {
  let result = page.evaluate(() => {
    return document.querySelector('#rightContent .pagination a.next');
  });

  return result === null ? null : result.href;
}

function extractProducts(page) {
  let results = page.evaluate(() => {
    return Array.from(
        document.querySelectorAll('.item_produits_courses > ul > li > a'),
    );
  });

  return results.map(r => r.href);
}

function extractProduct(page) {
  return page.evaluate(() => {
    let result = document.querySelector(
        '.contentProduct .description_produits.courses',
    );

    let product = {};

    let tmp = result.querySelector('h2.brand');
    if (tmp != null) {
      product.brand = tmp.innerText;
    }

    tmp = result.querySelector('h1');
    if (tmp != null) {
      product.title = tmp.innerText;
    }

    tmp = result.querySelector('h4');
    if (tmp != null) {
      product.subTitle = tmp.innerText;
    }

    tmp = result.querySelector('.description.courses #priceChange');
    if (tmp != null) {
      product.price = tmp.innerText;
    }

    tmp = result.querySelector('.weight-price');
    if (tmp != null) {
      product.weightPrice = tmp.innerText;
    }

    return product;
  });
}
