package com.fluxninja.example.controller;

import com.fluxninja.example.filter.ApertureFilter;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.core.env.Environment;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class AppController {
    @RequestMapping(value = "/super", method = RequestMethod.GET)
    public String hello() {
        return "Hello World";
    }

    @RequestMapping(value = "/health", method = RequestMethod.GET)
    public String health() {
        return "Healthy";
    }

    @RequestMapping(value = "/connected", method = RequestMethod.GET)
    public String connected() {
        return "";
    }

    @Bean
    public FilterRegistrationBean<ApertureFilter> apertureFilter(Environment env){
        FilterRegistrationBean<ApertureFilter> registrationBean = new FilterRegistrationBean<>();

        ApertureFilter filter = new ApertureFilter();
        registrationBean.setFilter(new ApertureFilter());
        registrationBean.addUrlPatterns("/super");

        String agentHost = env.getProperty("FN_AGENT_HOST");
        String agentPort = env.getProperty("FN_AGENT_PORT");
        registrationBean.addInitParameter("agent_host", agentHost);
        registrationBean.addInitParameter("agent_port", agentPort);

        return registrationBean;
    }
}
