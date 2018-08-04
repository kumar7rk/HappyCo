const puppeteer = require('puppeteer');
const CREDS = require('./creds');

(async () => {
	const browser = await puppeteer.launch({
		headless: false
	});

	const page = await browser.newPage();
	
	await page.goto("https://www.linkedin.com");

	const USERNAME_SELECTOR = '#login-email';
  	const PASSWORD_SELECTOR = '#login-password';
  	const BUTTON_SELECTOR = '#login-submit';

  	await page.click(USERNAME_SELECTOR);
  	await page.keyboard.type(CREDS.username);

  	await page.click(PASSWORD_SELECTOR);
  	await page.keyboard.type(CREDS.password);
	
	await page.click(BUTTON_SELECTOR);
    await page.waitForNavigation();
  	
	var url = "https://www.linkedin.com/in/robert-wise/";

	await page.goto(url);
	
	await page.waitFor(2 * 1000);
	await page.evaluate(_ => {
	  window.scrollBy(0, window.innerHeight);
	});
	await page.waitFor(2 * 1000);
   
   const title = await page.evaluate(() => document.querySelector
   	('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)
   
   console.log(title);

 await browser.close();
})();