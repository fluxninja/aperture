package com.fluxninja.aperture.instrumentation.armeria;

import static net.bytebuddy.matcher.ElementMatchers.*;

import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;
import net.bytebuddy.agent.builder.AgentBuilder;
import net.bytebuddy.description.type.TypeDescription;
import net.bytebuddy.matcher.ElementMatcher;
import net.bytebuddy.matcher.ElementMatchers;

public class ArmeriaServerInstrumentation implements TransformerInstrumentation {
    @Override
    public ElementMatcher<TypeDescription> getType() {
        return ElementMatchers.named("com.linecorp.armeria.server.ServerBuilder");
    }

    @Override
    public AgentBuilder.Transformer getTransformer() {
        return new AgentBuilder.Transformer.ForAdvice()
                .advice(
                        isMethod().and(isPublic()).and(named("build")),
                        ArmeriaServerAdvice.class.getName());
    }
}
