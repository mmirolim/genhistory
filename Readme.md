Generate sample data of users clicks from emails

It generates N number of profiles with int ids.
96% of profiles registration date is same and before calculation start date.
4% of profiles registration date is distributed normally during start week.
Registration dates generated randomly for each week with weibull distribution.
Every profile have Edg email domain one of yahoo, aol, hotmail,google.
Every click has coeff 1 and 0.2. Direct click to adv count as 1 others as 0.2.
Clicks time resolution is day to make it easier to calculate.
Clicks done every day by every profile number of clicks per day defined by weibull distribution.
Clicks and Unsub rate taken from http://mailchimp.com/resources/research/email-marketing-benchmarks/
If it game category click rate is CR 3.38% and unsub rate is UR 0.20%.
First on each day randomly UR number of profiles unsubscribe they are deleted from pool.
Then CR number taken randomly from the profile pools. They click according to some distribution.
Then new users registered by some distribution the rate is around 0.5%.

