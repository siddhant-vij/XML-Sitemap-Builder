# XML Sitemap Builder

[Gophercises](https://gophercises.com/) Exercise Details:

A sitemap is basically a map of all of the pages within a specific domain. They are used by search engines and other tools to inform them of all of the pages on your domain.

One way these can be built is by first visiting the root page of the website and making a list of every link on that page that goes to a page on the same domain.

Once you have created the list of links, you could then visit each and add any new links to your list. By repeating this step over and over (recursively) you would eventually visit every page that on the domain that can be reached by following links from the root page.

In this exercise your goal is to build a sitemap builder like the one described above. The end user will run the program and provide you with a URL (*hint - use a flag or a command line arg for this!*) that you will use to start the process.

Once you have determined all of the pages of a site, your sitemap builder should then output the data in the following XML format:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>http://www.example.com/dogs</loc>
    <priority>0.5</priority>
  </url>
</urlset>
```

*Note: This should be the same as the [standard sitemap protocol](https://www.sitemaps.org/index.html)*

Where each page is listed in its own `<url>` tag and includes the `<loc>` tag inside of it.

`<priority>` is an optional tag. The priority of a URL is relative to other URLs on your site. This value does not affect how your pages are compared to pages on other sites - it only lets the search engines know which pages you deem most important for the crawlers. Since the priority is relative, it is only used to select between URLs on your site.

*Come up with an algorithm to assign a priority between 0.0 and 1.0 - for a URL in the sitemap on the basis of the count of internal links to that page.*

In order to complete this exercise, use the [HTML Link Parser](https://github.com/siddhant-vij/HTML-Link-Parser) package created to parse your HTML pages for links.

<br>

## Technical Notes

- Be aware that links can be cyclical. That is, page `abc.com` may link to page `abc.com/about`, and then the about page may link back to the home page (`abc.com`). These cycles can also occur over many pages, for instance you might have:
  ```
  /about -> /contact
  /contact -> /pricing
  /pricing -> /testimonials
  /testimonials -> /about
  ```
  Where the cycle takes 4 links to finally reach it, but there is indeed a cycle. This is important to remember because you don't want your program to get into an infinite loop where it keeps visiting the same few pages over and over.
- The following packages will be helpful...
  - [net/http](https://golang.org/pkg/net/http/) - to initiate GET requests to each page in your sitemap and get the HTML on that page
  - [encoding/xml](https://golang.org/pkg/encoding/xml/) - to print out the XML output at the end
  - [flag](https://golang.org/pkg/flag/) - to parse user provided flags like the website domain

<br>

## Improvement Ideas

- In case of a very big website, sequentially building the sitemap can be very slow. Improve your sitemap builder to build the sitemap in parallel using goroutines & a concurrent data structure to handle the sitemap.
- Add in a `depth` flag that defines the maximum number of links to follow when building a sitemap. For instance, if you had a max depth of 3 and the following links:
  ```
  a->b->c->d
  ```

  Then your sitemap builder would not visit or include `d` because you must follow more than 3 links to to get to the page.

  On the other hand, if the links for the page were like this:

  ```
  a->b->c->d
  b->d
  ```

  Where there is also a link to page `d` from page `b`, then your sitemap builder should include `d` because it can be reached in 3 links.