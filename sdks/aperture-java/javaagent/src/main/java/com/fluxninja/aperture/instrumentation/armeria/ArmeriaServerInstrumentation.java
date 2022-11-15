package com.fluxninja.aperture.instrumentation.armeria;


import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.instrumentation.InstrumentationModule;
import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.ServerBuilder;
import net.bytebuddy.agent.builder.AgentBuilder;
import net.bytebuddy.asm.Advice;
import net.bytebuddy.description.NamedElement;
import net.bytebuddy.description.type.TypeDescription;
import net.bytebuddy.matcher.ElementMatcher;
import net.bytebuddy.matcher.ElementMatchers;

import static net.bytebuddy.matcher.ElementMatchers.*;


public class ArmeriaServerInstrumentation implements TransformerInstrumentation {
    @Override
    public ElementMatcher<TypeDescription> getType() {
        return ElementMatchers.named("com.linecorp.armeria.server.ServerBuilder");
    }

    @Override
    public AgentBuilder.Transformer getTransformer() {
        return new AgentBuilder.Transformer.ForAdvice()
                .advice(isMethod().and(isPublic()).and(named("build")),
                        ArmeriaServerAdvice.class.getName());
    }

}
